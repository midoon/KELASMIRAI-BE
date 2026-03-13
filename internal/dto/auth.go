package dto

import "mime/multipart"

type TenantRegisterRequest struct {
	SchoolName    string                `form:"school_name" validate:"required"`
	SchoolSlug    string                `form:"school_slug" validate:"required"`
	SchoolEmail   string                `form:"school_email" validate:"required,email"`
	SchoolPhone   string                `form:"school_phone"`
	SchoolAddress string                `form:"school_address"`
	SchoolLogo    *multipart.FileHeader `form:"school_logo"`

	AdminName     string `form:"admin_name" validate:"required"`
	AdminEmail    string `form:"admin_email" validate:"required,email"`
	AdminPassword string `form:"admin_password" validate:"required,min=8"`
}
