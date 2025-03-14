/*
** Copyright (c) 2022 Oracle and/or its affiliates.
**
** The Universal Permissive License (UPL), Version 1.0
**
** Subject to the condition set forth below, permission is hereby granted to any
** person obtaining a copy of this software, associated documentation and/or data
** (collectively the "Software"), free of charge and under any and all copyright
** rights in the Software, and any and all patent rights owned or freely
** licensable by each licensor hereunder covering either (i) the unmodified
** Software as contributed to or provided by such licensor, or (ii) the Larger
** Works (as defined below), to deal in both
**
** (a) the Software, and
** (b) any piece of software and/or hardware listed in the lrgrwrks.txt file if
** one is included with the Software (each a "Larger Work" to which the Software
** is contributed by such licensors),
**
** without restriction, including without limitation the rights to copy, create
** derivative works of, display, perform, and distribute the Software and make,
** use, sell, offer for sale, import, export, have made, and have sold the
** Software and the Larger Work(s), and to sublicense the foregoing rights on
** either these or other terms.
**
** This license is subject to the following condition:
** The above copyright notice and either this complete permission notice or at
** a minimum a reference to the UPL must be included in all copies or
** substantial portions of the Software.
**
** THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
** IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
** FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
** AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
** LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
** OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
** SOFTWARE.
 */

package e2etest

import (
	"context"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"
	"github.com/oracle/oci-go-sdk/v65/workrequests"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	dbv1alpha1 "github.com/oracle/oracle-database-operator/apis/database/v1alpha1"
	e2ebehavior "github.com/oracle/oracle-database-operator/test/e2e/behavior"
	e2eutil "github.com/oracle/oracle-database-operator/test/e2e/util"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var _ = Describe("test ADB binding with hardLink=true", func() {
	var adbLookupKey types.NamespacedName
	const downloadedWallet = "instance-wallet-secret-1"
	var adbID *string
	var terminatedAdbID string

	AfterEach(func() {
		// IMPORTANT: The operator might have to call reconcile multiple times to finish an operation.
		// If we do the update immediately, the previous reconciliation will overwrite the changes.
		By("Sleeping 20 seconds to wait for reconciliation to finish")
		time.Sleep(time.Second * 20)
	})

	It("should init the test", func() {
		By("creating a temp ADB in OCI for binding test")
		dbName := e2eutil.GenerateDBName()
		createResp, err := e2eutil.CreateAutonomousDatabase(dbClient, &SharedCompartmentOCID, &dbName, &SharedPlainTextAdminPassword)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(createResp.AutonomousDatabase.Id).ShouldNot(BeNil())

		By("Save the database ID for later use")
		adbID = createResp.AutonomousDatabase.Id
		terminatedAdbID = *adbID

		By("Wait until the work request is in SUCCEEDED status")
		workClient, err := workrequests.NewWorkRequestClientWithConfigurationProvider(configProvider)
		Expect(err).ShouldNot(HaveOccurred())

		err = e2eutil.WaitUntilWorkCompleted(workClient, createResp.OpcWorkRequestId)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("ADB binding with HardLink = false using Wallet Password Secret", func() {
		It("Should create a AutonomousDatabase resource with HardLink = false", func() {
			adb := &dbv1alpha1.AutonomousDatabase{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "database.oracle.com/v1alpha1",
					Kind:       "AutonomousDatabase",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bindadb",
					Namespace: ADBNamespace,
				},
				Spec: dbv1alpha1.AutonomousDatabaseSpec{
					Details: dbv1alpha1.AutonomousDatabaseDetails{
						Id: adbID,
					},
					Wallet: dbv1alpha1.WalletSpec{
						Name: common.String(downloadedWallet),
						Password: dbv1alpha1.PasswordSpec{
							K8sSecret: dbv1alpha1.K8sSecretSpec{
								Name: common.String(SharedWalletPassSecretName),
							},
						},
					},
					HardLink: common.Bool(false),
					OciConfig: dbv1alpha1.OciConfigSpec{
						ConfigMapName: common.String(SharedOCIConfigMapName),
						SecretName:    common.String(SharedOCISecretName),
					},
				},
			}

			adbLookupKey = types.NamespacedName{Name: adb.Name, Namespace: adb.Namespace}

			Expect(k8sClient.Create(context.TODO(), adb)).Should(Succeed())
		})

		It("should bind to an ADB", e2ebehavior.AssertBind(&k8sClient, &adbLookupKey))

		It("Should download an instance wallet using the password from K8s Secret "+SharedWalletPassSecretName, e2ebehavior.AssertWallet(&k8sClient, &adbLookupKey))

		It("should update ADB", e2ebehavior.UpdateAndAssertDetails(&k8sClient, &dbClient, &adbLookupKey, SharedNewAdminPassSecretName, &SharedPlainTextNewAdminPassword, &SharedPlainTextWalletPassword))

		It("Should stop ADB", e2ebehavior.UpdateAndAssertADBState(&k8sClient, &dbClient, &adbLookupKey, database.AutonomousDatabaseLifecycleStateStopped))

		It("Should restart ADB", e2ebehavior.UpdateAndAssertADBState(&k8sClient, &dbClient, &adbLookupKey, database.AutonomousDatabaseLifecycleStateAvailable))

		It("Should change to RESTRICTED network access", e2ebehavior.TestNetworkAccessRestricted(&k8sClient, &dbClient, &adbLookupKey, []string{"192.168.0.1"}, false))

		It("Should change isMTLSConnectionRequired to false", e2ebehavior.TestNetworkAccessRestricted(&k8sClient, &dbClient, &adbLookupKey, []string{"192.168.0.1"}, false))

		It("Should should change to PRIVATE network access", e2ebehavior.TestNetworkAccessPrivate(&k8sClient, &dbClient, &adbLookupKey, false, &SharedSubnetOCID, &SharedNsgOCID))

		It("Should change isMTLSConnectionRequired to true when network access is PRIVATE", e2ebehavior.TestNetworkAccessPrivate(&k8sClient, &dbClient, &adbLookupKey, true, &SharedSubnetOCID, &SharedNsgOCID))

		It("Should return to PUBLIC access type", e2ebehavior.TestNetworkAccessPublic(&k8sClient, &dbClient, &adbLookupKey))

		It("Should delete the resource in cluster but not terminate the database in OCI", e2ebehavior.AssertSoftLinkDelete(&k8sClient, &adbLookupKey))
	})

	Describe("ADB binding with HardLink = true using Wallet Password OCID", func() {
		It("Should create a AutonomousDatabase resource with HardLink = true", func() {
			adb := &dbv1alpha1.AutonomousDatabase{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "database.oracle.com/v1alpha1",
					Kind:       "AutonomousDatabase",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bindadb",
					Namespace: ADBNamespace,
				},
				Spec: dbv1alpha1.AutonomousDatabaseSpec{
					Details: dbv1alpha1.AutonomousDatabaseDetails{
						Id: adbID,
					},
					Wallet: dbv1alpha1.WalletSpec{
						Name: common.String(downloadedWallet),
						Password: dbv1alpha1.PasswordSpec{
							OciSecret: dbv1alpha1.OciSecretSpec{
								Id: common.String(SharedInstanceWalletPasswordOCID),
							},
						},
					},
					HardLink: common.Bool(true),
					OciConfig: dbv1alpha1.OciConfigSpec{
						ConfigMapName: common.String(SharedOCIConfigMapName),
						SecretName:    common.String(SharedOCISecretName),
					},
				},
			}

			adbLookupKey = types.NamespacedName{Name: adb.Name, Namespace: adb.Namespace}

			Expect(k8sClient.Create(context.TODO(), adb)).Should(Succeed())
		})

		It("should bind to an ADB", e2ebehavior.AssertBind(&k8sClient, &adbLookupKey))

		It("Should download an instance wallet using the password from OCI Secret OCID "+SharedInstanceWalletPasswordOCID, e2ebehavior.AssertWallet(&k8sClient, &adbLookupKey))

		It("Should delete the resource in cluster and terminate the database in OCI", e2ebehavior.AssertHardLinkDelete(&k8sClient, &dbClient, &adbLookupKey))
	})

	//Bind to terminated adb from previous test
	Describe("bind to a terminated adb", func() {

		//Wait until remote state is terminated
		It("Should check that OCI adb state is terminated", e2ebehavior.AssertADBRemoteStateOCID(&k8sClient, &dbClient, &terminatedAdbID, database.AutonomousDatabaseLifecycleStateTerminated, time.Second*300))

		It("Should create a AutonomousDatabase resource", func() {
			adb := &dbv1alpha1.AutonomousDatabase{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "database.oracle.com/v1alpha1",
					Kind:       "AutonomousDatabase",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bindadb",
					Namespace: ADBNamespace,
				},
				Spec: dbv1alpha1.AutonomousDatabaseSpec{
					Details: dbv1alpha1.AutonomousDatabaseDetails{
						Id: &terminatedAdbID,
					},
					HardLink: common.Bool(true),
					OciConfig: dbv1alpha1.OciConfigSpec{
						ConfigMapName: common.String(SharedOCIConfigMapName),
						SecretName:    common.String(SharedOCISecretName),
					},
				},
			}

			adbLookupKey = types.NamespacedName{Name: adb.Name, Namespace: adb.Namespace}

			Expect(k8sClient.Create(context.TODO(), adb)).Should(Succeed())
		})

		It("Should check for TERMINATED state in local resource", e2ebehavior.AssertADBLocalState(&k8sClient, &adbLookupKey, database.AutonomousDatabaseLifecycleStateTerminated))

		It("Should delete local resource", e2ebehavior.AssertSoftLinkDelete(&k8sClient, &adbLookupKey))
	})
})
