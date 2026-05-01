package externalcontact

const (
	EventChangeExternalContact = "change_external_contact"
	EventChangeExternalChat    = "change_external_chat"

	ChangeTypeAddExternalContact     = "add_external_contact"
	ChangeTypeEditExternalContact    = "edit_external_contact"
	ChangeTypeAddHalfExternalContact = "add_half_external_contact"
	ChangeTypeDelExternalContact     = "del_external_contact"
	ChangeTypeDelFollowUser          = "del_follow_user"
	ChangeTypeTransferFail           = "transfer_fail"
	ChangeTypeCreate                 = "create"
	ChangeTypeUpdate                 = "update"
	ChangeTypeDismiss                = "dismiss"
	ChangeTypeDelete                 = "delete"
	ChangeTypeShuffle                = "shuffle"

	UpdateDetailAddMember    = "add_member"
	UpdateDetailDelMember    = "del_member"
	UpdateDetailChangeOwner  = "change_owner"
	UpdateDetailChangeName   = "change_name"
	UpdateDetailChangeNotice = "change_notice"
)
