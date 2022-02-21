/*
** Copyright (c) 2021 Oracle and/or its affiliates.
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

package v1alpha1

import (
	"encoding/json"
	"strconv"

	"github.com/oracle/oci-go-sdk/v54/database"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/oracle/oracle-database-operator/commons/annotations"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AutonomousDatabaseSpec defines the desired state of AutonomousDatabase
// Important: Run "make" to regenerate code after modifying this file
type AutonomousDatabaseSpec struct {
	Details   AutonomousDatabaseDetails `json:"details"`
	OCIConfig OCIConfigSpec             `json:"ociConfig,omitempty"`
	// +kubebuilder:default:=false
	HardLink *bool `json:"hardLink,omitempty"`
}

type OCIConfigSpec struct {
	ConfigMapName *string `json:"configMapName,omitempty"`
	SecretName    *string `json:"secretName,omitempty"`
}

// AutonomousDatabaseDetails defines the detail information of AutonomousDatabase, corresponding to oci-go-sdk/database/AutonomousDatabase
type AutonomousDatabaseDetails struct {
	AutonomousDatabaseOCID          *string `json:"autonomousDatabaseOCID,omitempty"`
	CompartmentOCID                 *string `json:"compartmentOCID,omitempty"`
	AutonomousContainerDatabaseOCID *string `json:"autonomousContainerDatabaseOCID,omitempty"`
	DisplayName                     *string `json:"displayName,omitempty"`
	DbName                          *string `json:"dbName,omitempty"`
	// +kubebuilder:validation:Enum:="OLTP";"DW";"AJD";"APEX"
	DbWorkload           database.AutonomousDatabaseDbWorkloadEnum     `json:"dbWorkload,omitempty"`
	IsDedicated          *bool                                         `json:"isDedicated,omitempty"`
	DbVersion            *string                                       `json:"dbVersion,omitempty"`
	DataStorageSizeInTBs *int                                          `json:"dataStorageSizeInTBs,omitempty"`
	CPUCoreCount         *int                                          `json:"cpuCoreCount,omitempty"`
	AdminPassword        PasswordSpec                                  `json:"adminPassword,omitempty"`
	IsAutoScalingEnabled *bool                                         `json:"isAutoScalingEnabled,omitempty"`
	LifecycleState       database.AutonomousDatabaseLifecycleStateEnum `json:"lifecycleState,omitempty"`

	NetworkAccess NetworkAccessSpec `json:"networkAccess,omitempty"`

	FreeformTags map[string]string `json:"freeformTags,omitempty"`

	Wallet WalletSpec `json:"wallet,omitempty"`
}

type WalletSpec struct {
	Name     *string      `json:"name,omitempty"`
	Password PasswordSpec `json:"password,omitempty"`
}

type PasswordSpec struct {
	K8sSecretName *string `json:"k8sSecretName,omitempty"`
	OCISecretOCID *string `json:"ociSecretOCID,omitempty"`
}

type NetworkAccessTypeEnum string

const (
	NetworkAccessTypePublic     NetworkAccessTypeEnum = "PUBLIC"
	NetworkAccessTypeRestricted NetworkAccessTypeEnum = "RESTRICTED"
	NetworkAccessTypePrivate    NetworkAccessTypeEnum = "PRIVATE"
)

type NetworkAccessSpec struct {
	// +kubebuilder:validation:Enum:="";"PUBLIC";"RESTRICTED";"PRIVATE"
	AccessType               NetworkAccessTypeEnum `json:"accessType,omitempty"`
	IsAccessControlEnabled   *bool                 `json:"isAccessControlEnabled,omitempty"`
	AccessControlList        []string              `json:"accessControlList,omitempty"`
	PrivateEndpoint          PrivateEndpointSpec   `json:"privateEndpoint,omitempty"`
	IsMTLSConnectionRequired *bool                 `json:"isMTLSConnectionRequired,omitempty"`
}

type PrivateEndpointSpec struct {
	SubnetOCID     *string  `json:"subnetOCID,omitempty"`
	NsgOCIDs       []string `json:"nsgOCIDs,omitempty"`
	HostnamePrefix *string  `json:"hostnamePrefix,omitempty"`
}

// AutonomousDatabaseStatus defines the observed state of AutonomousDatabase
type AutonomousDatabaseStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DisplayName          string                                        `json:"displayName,omitempty"`
	LifecycleState       database.AutonomousDatabaseLifecycleStateEnum `json:"lifecycleState,omitempty"`
	IsDedicated          string                                        `json:"isDedicated,omitempty"`
	CPUCoreCount         int                                           `json:"cpuCoreCount,omitempty"`
	DataStorageSizeInTBs int                                           `json:"dataStorageSizeInTBs,omitempty"`
	DbWorkload           database.AutonomousDatabaseDbWorkloadEnum     `json:"dbWorkload,omitempty"`
	TimeCreated          string                                        `json:"timeCreated,omitempty"`
	AllConnectionStrings []ConnectionStringsSet                        `json:"allConnectionStrings,omitempty"`
}

type TLSAuthenticationEnum string

const (
	TLSAuthenticationTLS  TLSAuthenticationEnum = "TLS"
	TLSAuthenticationmTLS TLSAuthenticationEnum = "Mutual TLS"
)

type ConnectionStringsSet struct {
	TLSAuthentication TLSAuthenticationEnum `json:"tlsAuthentication,omitempty"`
	ConnectionStrings map[string]string     `json:"connectionStrings"`
}

// AutonomousDatabase is the Schema for the autonomousdatabases API
// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName="adb";"adbs"
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".status.displayName",name="Display Name",type=string
// +kubebuilder:printcolumn:JSONPath=".status.lifecycleState",name="State",type=string
// +kubebuilder:printcolumn:JSONPath=".status.isDedicated",name="Dedicated",type=string
// +kubebuilder:printcolumn:JSONPath=".status.cpuCoreCount",name="OCPUs",type=integer
// +kubebuilder:printcolumn:JSONPath=".status.dataStorageSizeInTBs",name="Storage (TB)",type=integer
// +kubebuilder:printcolumn:JSONPath=".status.dbWorkload",name="Workload Type",type=string
// +kubebuilder:printcolumn:JSONPath=".status.timeCreated",name="Created",type=string
type AutonomousDatabase struct {
	metaV1.TypeMeta   `json:",inline"`
	metaV1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AutonomousDatabaseSpec   `json:"spec,omitempty"`
	Status AutonomousDatabaseStatus `json:"status,omitempty"`
}

// LastSuccessfulSpec is an annotation key which maps to the value of last successful spec
const LastSuccessfulSpec string = "lastSuccessfulSpec"

// GetLastSuccessfulSpec returns spec from the lass successful reconciliation.
// Returns nil, nil if there is no lastSuccessfulSpec.
func (adb *AutonomousDatabase) GetLastSuccessfulSpec() (*AutonomousDatabaseSpec, error) {
	val, ok := adb.GetAnnotations()[LastSuccessfulSpec]
	if !ok {
		return nil, nil
	}

	specBytes := []byte(val)
	sucSpec := AutonomousDatabaseSpec{}

	err := json.Unmarshal(specBytes, &sucSpec)
	if err != nil {
		return nil, err
	}

	return &sucSpec, nil
}

// UpdateLastSuccessfulSpec updates lastSuccessfulSpec with the current spec.
func (adb *AutonomousDatabase) UpdateLastSuccessfulSpec(kubeClient client.Client) error {
	specBytes, err := json.Marshal(adb.Spec)
	if err != nil {
		return err
	}

	anns := map[string]string{
		LastSuccessfulSpec: string(specBytes),
	}

	return annotations.SetAnnotations(kubeClient, adb, anns)
}

// UpdateAttrFromOCIAutonomousDatabase updates the attributes from database.AutonomousDatabase object and returns the resource
func (adb *AutonomousDatabase) UpdateAttrFromOCIAutonomousDatabase(ociObj database.AutonomousDatabase) *AutonomousDatabase {
	/***********************************
	* update the spec
	***********************************/
	adb.Spec.Details.AutonomousDatabaseOCID = ociObj.Id
	adb.Spec.Details.CompartmentOCID = ociObj.CompartmentId
	adb.Spec.Details.AutonomousContainerDatabaseOCID = ociObj.AutonomousContainerDatabaseId
	adb.Spec.Details.DisplayName = ociObj.DisplayName
	adb.Spec.Details.DbName = ociObj.DbName
	adb.Spec.Details.DbWorkload = ociObj.DbWorkload
	adb.Spec.Details.IsDedicated = ociObj.IsDedicated
	adb.Spec.Details.DbVersion = ociObj.DbVersion
	adb.Spec.Details.DataStorageSizeInTBs = ociObj.DataStorageSizeInTBs
	adb.Spec.Details.CPUCoreCount = ociObj.CpuCoreCount
	adb.Spec.Details.IsAutoScalingEnabled = ociObj.IsAutoScalingEnabled
	adb.Spec.Details.LifecycleState = ociObj.LifecycleState
	adb.Spec.Details.FreeformTags = ociObj.FreeformTags

	if *ociObj.IsDedicated {
		adb.Spec.Details.NetworkAccess.AccessType = NetworkAccessTypePrivate
	} else {
		if ociObj.NsgIds != nil {
			adb.Spec.Details.NetworkAccess.AccessType = NetworkAccessTypePrivate
		} else if ociObj.WhitelistedIps != nil {
			adb.Spec.Details.NetworkAccess.AccessType = NetworkAccessTypeRestricted
		} else {
			adb.Spec.Details.NetworkAccess.AccessType = NetworkAccessTypePublic
		}
	}

	adb.Spec.Details.NetworkAccess.IsAccessControlEnabled = ociObj.IsAccessControlEnabled
	adb.Spec.Details.NetworkAccess.AccessControlList = ociObj.WhitelistedIps
	adb.Spec.Details.NetworkAccess.IsMTLSConnectionRequired = ociObj.IsMtlsConnectionRequired
	adb.Spec.Details.NetworkAccess.PrivateEndpoint.SubnetOCID = ociObj.SubnetId
	adb.Spec.Details.NetworkAccess.PrivateEndpoint.NsgOCIDs = ociObj.NsgIds
	adb.Spec.Details.NetworkAccess.PrivateEndpoint.HostnamePrefix = ociObj.PrivateEndpointLabel

	/***********************************
	* update the status subresource
	***********************************/
	adb.Status.DisplayName = *ociObj.DisplayName
	adb.Status.LifecycleState = ociObj.LifecycleState
	adb.Status.IsDedicated = strconv.FormatBool(*ociObj.IsDedicated)
	adb.Status.CPUCoreCount = *ociObj.CpuCoreCount
	adb.Status.DataStorageSizeInTBs = *ociObj.DataStorageSizeInTBs
	adb.Status.DbWorkload = ociObj.DbWorkload
	adb.Status.TimeCreated = ociObj.TimeCreated.String()

	var curAlllConns []ConnectionStringsSet
	if *ociObj.IsDedicated {
		connSet := ConnectionStringsSet{ConnectionStrings: ociObj.ConnectionStrings.AllConnectionStrings}
		curAlllConns = append(curAlllConns, connSet)

	} else {
		mTLSStrings := make(map[string]string)
		tlsStrings := make(map[string]string)

		for _, profile := range ociObj.ConnectionStrings.Profiles {
			if profile.TlsAuthentication == database.DatabaseConnectionStringProfileTlsAuthenticationMutual {
				mTLSStrings[*profile.DisplayName] = *profile.Value
			} else {
				tlsStrings[*profile.DisplayName] = *profile.Value
			}
		}

		if len(mTLSStrings) > 0 {
			mTLSConnSet := ConnectionStringsSet{
				TLSAuthentication: TLSAuthenticationmTLS,
				ConnectionStrings: mTLSStrings,
			}

			curAlllConns = append(curAlllConns, mTLSConnSet)
		}

		if len(tlsStrings) > 0 {
			tlsConnSet := ConnectionStringsSet{
				TLSAuthentication: TLSAuthenticationTLS,
				ConnectionStrings: tlsStrings,
			}

			curAlllConns = append(curAlllConns, tlsConnSet)
		}
	}
	adb.Status.AllConnectionStrings = curAlllConns

	return adb
}

// +kubebuilder:object:root=true

// AutonomousDatabaseList contains a list of AutonomousDatabase
type AutonomousDatabaseList struct {
	metaV1.TypeMeta `json:",inline"`
	metaV1.ListMeta `json:"metadata,omitempty"`
	Items           []AutonomousDatabase `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AutonomousDatabase{}, &AutonomousDatabaseList{})
}

// A helper function which is useful for debugging. The function prints out a structural JSON format.
func (adb *AutonomousDatabase) String() (string, error) {
	out, err := json.MarshalIndent(adb, "", "    ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}
