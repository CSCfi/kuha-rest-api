package authz

var RolePermissions = map[string][]string{
	// Admin
	"admin": {"*"},

	// FIS roles
	"fis": {
		"GET:/v1/fis",
		"POST:/v1/fis",
		"PUT:/v1/fis",
		"DELETE:/v1/fis",
	},
	"fis_read": {
		"GET:/v1/fis",
	},

	// UTV roles
	"utv": {
		"GET:/v1/utv",
		"POST:/v1/utv",
		"PUT:/v1/utv",
		"DELETE:/v1/utv",
	},
	"utv_read": {
		"GET:/v1/utv",
	},

	// KAMK roles
	"kamk": {
		"GET:/v1/kamk",
		"POST:/v1/kamk",
		"PUT:/v1/kamk",
		"DELETE:/v1/kamk",
	},
	"kamk_read": {
		"GET:/v1/kamk",
	},

	// K-LAB roles
	"klab": {
		"GET:/v1/klab",
		"POST:/v1/klab",
		"PUT:/v1/klab",
		"DELETE:/v1/klab",
	},
	"klab_read": {
		"GET:/v1/klab",
	},

	// M360 roles
	"m360": {
		"GET:/v1/m360",
		"POST:/v1/m360",
		"PUT:/v1/m360",
		"DELETE:/v1/m360",
	},
	"m360_read": {
		"GET:/v1/m360",
	},

	// Couchtech roles
	"couchtech": {
		"GET:/v1/couchtech",
		"POST:/v1/couchtech",
		"PUT:/v1/couchtech",
		"DELETE:/v1/couchtech",
	},
	"couchtech_read": {
		"GET:/v1/couchtech",
	},

	// Archinisis roles
	"archinisis": {
		"GET:/v1/archinisis",
		"POST:/v1/archinisis",
		"PUT:/v1/archinisis",
		"DELETE:/v1/archinisis",
	},
	"archinisis_read": {
		"GET:/v1/archinisis",
	},
}
