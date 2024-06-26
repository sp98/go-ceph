//go:build !nautilus
// +build !nautilus

package rbd

// #cgo LDFLAGS: -lrbd
// #include <stdlib.h>
// #include <rbd/librbd.h>
import "C"
import (
	"github.com/ceph/go-ceph/rados"
)

// MirrorGroupStatusState is used to indicate the state of a mirrored group
// within the site status info.
type MirrorGroupStatusState int64

const (
	// MirrorGrouptatusStateUnknown is equivalent to MIRROR_GROUP_STATUS_STATE_UNKNOWN
	MirrorGroupStatusStateUnknown = MirrorImageStatusState(C.MIRROR_GROUP_STATUS_STATE_UNKNOWN)
	// MirrorGroupStatusStateError is equivalent to MIRROR_GROUP_STATUS_STATE_ERROR
	MirrorGroupStatusStateError = MirrorImageStatusState(C.MIRROR_GROUP_STATUS_STATE_ERROR)
	// MirrorGroupStatusStateStartingReplay is equivalent to MIRROR_GROUP_STATUS_STATE_STARTING_REPLAY
	MirrorGroupStatusStateStartingReplay = MirrorImageStatusState(C.MIRROR_GROUP_STATUS_STATE_STARTING_REPLAY)
	// MirrorGroupStatusStateReplaying is equivalent to MIRROR_GROUP_STATUS_STATE_REPLAYING
	MirrorGroupStatusStateReplaying = MirrorImageStatusState(C.MIRROR_GROUP_STATUS_STATE_REPLAYING)
	// MirrorGroupStatusStateStoppingReplay is equivalent to MIRROR_GROUP_STATUS_STATE_STOPPING_REPLAY
	MirrorGroupStatusStateStoppingReplay = MirrorImageStatusState(C.MIRROR_GROUP_STATUS_STATE_STOPPING_REPLAY)
	// MirrorGroupStatusStateStopped is equivalent to MIRROR_IMAGE_GROUP_STATUS_STATE_STOPPED
	MirrorGroupStatusStateStopped = MirrorImageStatusState(C.MIRROR_IMAGE_GROUP_STATUS_STATE_STOPPED)
)

// MirrorImageState represents the mirroring state of a RBD image.
type MirrorGroupState C.rbd_mirror_group_state_t

const (
	// MirrorGrpupDisabling is the representation of
	// RBD_MIRROR_GROUP_DISABLING from librbd.
	MirrorGrpupDisabling = MirrorGroupState(C.RBD_MIRROR_GROUP_DISABLING)
	// MirrorGroupEnabling is the representation of
	// RBD_MIRROR_GROUP_ENABLING from librbd
	MirrorGroupEnabling = MirrorGroupState(C.RBD_MIRROR_GROUP_ENABLING)
	// MirrorGroupEnabled is the representation of
	// RBD_MIRROR_IMAGE_ENABLED from librbd.
	MirrorGroupEnabled = MirrorGroupState(C.RBD_MIRROR_GROUP_ENABLED)
	// MirrorGroupDisabled is the representation of
	// RBD_MIRROR_GROUP_DISABLED from librbd.
	MirrorGroupDisabled = MirrorGroupState(C.RBD_MIRROR_GROUP_DISABLED)
)

// MirrorGroupInfo represents the mirroring status information of group.
type MirrorGroupInfo struct {
	GlobalID string
	State    MirrorGroupState
	Primary  bool
}

// SiteMirrorGroupStatus contains information pertaining to the status of
// a mirrored group within a site.
type SiteMirrorGroupStatus struct {
	MirrorUUID         string
	State              MirrorGroupStatusState
	MirrorImageCount   int
	MirrorImagePoolIds int64
	Description        string
	LastUpdate         int64
	Up                 bool
}

// GlobalMirrorGroupStatus contains information pertaining to the global
// status of a mirrored group. It contains general information as well
// as per-site information stored in the SiteStatuses slice.
type GlobalMirrorGroupStatus struct {
	Name              string
	Info              MirrorGroupInfo
	SiteStatusesCount int
	SiteStatuses      []SiteMirrorGroupStatus
}

// GetGlobalMirrorGroupStatus returns status information pertaining to the state
// of a groups's mirroring.
//
// Implements:
//
//	int rbd_mirror_group_get_status(
//		IoCtx& io_ctx,
//		const char *group_name
//		mirror_group_global_status_t *mirror_group_status,
//		size_t status_size);
func GetGlobalMirrorGroupStatus(ioctx *rados.IOContext, groupName string) (GlobalMirrorGroupStatus, error) {
	s := C.rbd_mirror_group_global_status_t{}
	ret := C.rbd_mirror_group_get_global_status(
		ioctx,
		groupName,
		&s,
		C.sizeof_rbd_mirror_group_global_status_t)
	if err := getError(ret); err != nil {
		return GlobalMirrorGroupStatus{}, err
	}

	// status := newGlobalMirrorGroupStatus(&s)
	return GlobalMirrorGroupStatus{}, nil
	// return status, nil
}

// func newGlobalMirrorGroupStatus(
// 	s *C.rbd_mirror_group_global_status_t) GlobalMirrorGroupStatus {

// 	status := GlobalMirrorGroupStatus{
// 		Name:         C.GoString(s.name),
// 		Info:         convertMirrorGroupInfo(&s.info),
// 		SiteStatuses: make([]SiteMirrorImageStatus, s.site_statuses_count),
// 	}
// 	// use the "Sven Technique" to treat the C pointer as a go slice temporarily
// 	sscs := (*siteArray)(unsafe.Pointer(s.site_statuses))[:s.site_statuses_count:s.site_statuses_count]
// 	for i := C.uint32_t(0); i < s.site_statuses_count; i++ {
// 		ss := sscs[i]
// 		status.SiteStatuses[i] = SiteMirrorImageStatus{
// 			MirrorUUID:  C.GoString(ss.mirror_uuid),
// 			State:       MirrorImageStatusState(ss.state),
// 			Description: C.GoString(ss.description),
// 			LastUpdate:  int64(ss.last_update),
// 			Up:          bool(ss.up),
// 		}
// 	}
// 	return status
// }

// func convertMirrorGroupInfo(cInfo *C.rbd_mirror_group_info_t) MirrorImageInfo {
// 	return MirrorGroupInfo{
// 		GlobalID: C.GoString(cInfo.global_id),
// 		State:    MirrorImageState(cInfo.state),
// 		Primary:  bool(cInfo.primary),
// 	}
// }
