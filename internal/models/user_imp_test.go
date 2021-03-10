package models

import "testing"

func TestUser_Validate(t *testing.T) {
	type fields struct {
		LoginID  string
		Password string
	}
	tests := []struct {
		fields  fields
		wantErr bool
	}{
		// valid login_id
		{fields{"lowercase_only", "OK_Password"}, false},
		{fields{"UPPERCASE_ONLY", "OK_Password"}, false},
		{fields{"Alphabet_and_0123456789", "OK_Password"}, false},
		{fields{"Contains_usable_char_-'.", "OK_Password"}, false},
		{fields{"ContainUppercaseAndLowercase", "OK_Password"}, false},
		{fields{"6xxxxX", "OK_Password"}, false},
		{fields{"30xxxxxxxxxxxxxxxxxxxxxxxxxxxX", "OK_Password"}, false},

		// invalid login_id
		{fields{"", "OK_Password"}, true},
		{fields{"5xxxX", "OK_Password"}, true},
		{fields{"31xxxxxxxxxxxxxxxxxxxxxxxxxxxxX", "OK_Password"}, true},
		{fields{"日本語", "OK_Password"}, true},

		// valid password
		{fields{"OK_user_name", "ContainUppercaseAndLowercase"}, false},
		{fields{"OK_user_name", "Alphabet_and_0123456789"}, false},
		{fields{"OK_user_name", "Contains_usable_char_-'."}, false},
		{fields{"8xxxxxxX", "OK_Password"}, false},
		{fields{"OK_user_name", "30xxxxxxxxxxxxxxxxxxxxxxxxxxxX"}, false},

		// invalid password
		{fields{"OK_user_name", ""}, true},
		{fields{"OK_user_name", "lowercase_only"}, true},
		{fields{"OK_user_name", "UPPERCASE_ONLY"}, true},
		{fields{"OK_user_name", "Contains_OK_user_name"}, true},
		{fields{"OK_user_name", "7xxxxxX"}, true},
		{fields{"OK_user_name", "31xxxxxxxxxxxxxxxxxxxxxxxxxxxxX"}, true},
		{fields{"OK_user_name", "日本語"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_\""}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_\n"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_!"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_@"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_#"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_%"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_$"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_^"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_&"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_*"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_("}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_)"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_+"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_,"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_;"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_:"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_:["}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_:]"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_:|"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_:?"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_:/"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_:<"}, true},
		{fields{"OK_user_name", "Contains_Unusable_char_:>"}, true},
	}
	for i, test := range tests {
		u := &User{
			LoginID:  test.fields.LoginID,
			Password: test.fields.Password,
		}
		if err := u.Validate(); (err != nil) != test.wantErr {
			t.Errorf("#%d Validate() error = %v, wantErr %v", i, err, test.wantErr)
		}
	}
}
