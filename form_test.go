package ditto_test

import (
	"encoding/json"
	"github.com/payfazz/ditto"
	"testing"
)

func TestForm(t *testing.T) {
	var structure map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &structure)
	if nil != err {
		t.Fatal(err)
	}

	f, err := ditto.NewSectionFromMap(structure)
	if nil != err {
		t.Fatal(err)
	}

	t.Logf("%+v", f)

	m, err := json.Marshal(f)
	if nil != err {
		t.Fatal(err)
	}
	t.Log(string(m))
}

var jsonData = `{
  "id": "task_form_section",
  "type": "summary_section_send",
  "child": [
    {
      "id": "agent_data_section",
      "info": {
        "icon": "https://content.payfazz.com/object/e1a14d808797311aaba5f164d50e30fb7ec037a15b07f6e687f31eaa49b4c60f"
      },
      "type": "nextable_form",
      "child": [
        {
          "id": "business_name",
          "type": "text",
          "title": "Nama Toko",
          "description": "Masukkan Nama Toko",
          "placeholder": "Nama Toko",
          "validations": [
            {
              "type": "required",
              "error_message": "Harus diisi"
            },
            {
              "type": "regex",
              "value": "%5E%5B%5Cw%20%5D%7B2%2C100%7D%24",
              "error_message": "Pastikan nama toko alphanumeric"
            }
          ]
        },
        {
          "id": "business_type",
          "info": {
            "hint": "",
            "icon": "https://content.payfazz.com/object/8a743b686f9edaeb5b309cb9318088f0bbbfb40ec7cdd9735ed40b8233528783",
            "options": {
              "type": "static",
              "value": "business_type"
            }
          },
          "type": "object_searchable_list",
          "title": "Jenis Toko",
          "description": "Jenis Toko",
          "placeholder": "Jenis Toko",
          "validations": [
            {
              "type": "required",
              "error_message": "Harus diisi"
            }
          ]
        },
        {
          "id": "name",
          "info": {
            "icon": "https://content.payfazz.com/object/8a743b686f9edaeb5b309cb9318088f0bbbfb40ec7cdd9735ed40b8233528783"
          },
          "type": "text",
          "title": "Nama Pemilik Toko",
          "description": "Masukkan Nama Pemilik Toko",
          "placeholder": "Nama Pemilik Toko",
          "validations": [
            {
              "type": "required",
              "error_message": "Harus diisi"
            },
            {
              "type": "regex",
              "value": "%5E%5BA-Za-z%20%5D%7B2%2C100%7D%24",
              "error_message": "Pastikan nama pemilik toko alphabet"
            }
          ]
        },
        {
          "id": "phone",
          "type": "text_numeric",
          "title": "Nomor Handphone Pemilik Toko",
          "description": "Masukkan Nomor Handphone Pemilik Toko",
          "placeholder": "081234567890",
          "validations": [
            {
              "type": "required",
              "error_message": "Harus diisi"
            },
            {
              "type": "regex",
              "value": "%5E%28%5C%2B%3F62%7C0%29%5B0-9%5D%7B10%2C13%7D%24",
              "error_message": "Nomor Telepon tidak valid"
            }
          ]
        },
        {
          "id": "location",
          "info": {
            "icon": "https://content.payfazz.com/object/8a743b686f9edaeb5b309cb9318088f0bbbfb40ec7cdd9735ed40b8233528783",
            "options": {
              "type": "dynamic",
              "value": "area"
            }
          },
          "type": "object_searchable_list",
          "title": "Daerah Tempat Toko Berada (Kelurahan)",
          "value": "",
          "description": "Masukkan daerah tempat toko berada (Kelurahan)",
          "placeholder": "Contoh: Gandaria Utara, Kebayoran baru, Jakarta Selatan, DKI Jakarta, Indonesia",
          "validations": [
            {
              "type": "required",
              "error_message": "Harus diisi"
            }
          ]
        },
        {
          "id": "address",
          "type": "text_multiline",
          "title": "Daerah Tempat Tinggal Agen",
          "description": "Masukkan Daerah Tempat Tinggal Agen",
          "placeholder": "Contoh: Jl, RT, Rw",
          "validations": [
            {
              "type": "required",
              "error_message": "Harus diisi"
            },
            {
              "type": "regex",
              "value": "%5E%5BA-Za-z0-9%5Cs.%5C%2C%2F-%5D%7B2%2C200%7D%24",
              "error_message": "Alamat tidak valid"
            },
            {
              "type": "count_between",
              "value": "2,200",
              "error_message": "Alamat tidak valid"
            }
          ]
        },
        {
          "id": "zip_code",
          "type": "text_numeric",
          "title": "Kode Pos",
          "description": "Masukkan Kode Pos",
          "placeholder": "00000",
          "validations": [
            {
              "type": "required",
              "error_message": "Harus diisi"
            },
            {
              "type": "regex",
              "value": "%5E%5Cd%2B%24",
              "error_message": "Input numeric aja"
            },
            {
              "type": "count_between",
              "value": "5,5",
              "error_message": "Harus 5 digit"
            }
          ]
        }
      ],
      "title": "Masukkan data Agen",
      "description": null
    },
    {
      "id": "take_photo_store_section",
      "info": {
        "icon": "https://content.payfazz.com/object/abbe92f80a439db37a1b81131feaa07c125b880af60b494cbffb02d79ecec807"
      },
      "type": "nextable_form",
      "child": [
        {
          "id": "business_image",
          "info": {
            "hint": "",
            "icon": "https://content.payfazz.com/object/abbe92f80a439db37a1b81131feaa07c125b880af60b494cbffb02d79ecec807",
            "instruction_text": "",
            "instruction_image": "https://content.payfazz.com/object/ac6575b50215816e6bf2054770969c1e9c27d0580d72761641332d679b82f51d"
          },
          "type": "photo_camera",
          "title": "Foto Toko",
          "description": "Ambil foto toko dengan jelas dan tidak terpotong",
          "placeholder": null,
          "validations": [
            {
              "type": "required",
              "error_message": "Harus diisi"
            }
          ]
        }
      ],
      "title": "Foto Tampak Depan Toko",
      "description": null
    },
    {
      "id": "task_store_owner_photo_section",
      "info": {
        "icon": "https://content.payfazz.com/object/abbe92f80a439db37a1b81131feaa07c125b880af60b494cbffb02d79ecec807"
      },
      "type": "nextable_form",
      "child": [
        {
          "id": "owner_image",
          "info": {
            "hint": "",
            "icon": "https://content.payfazz.com/object/abbe92f80a439db37a1b81131feaa07c125b880af60b494cbffb02d79ecec807",
            "instruction_text": "",
            "instruction_image": "https://content.payfazz.com/object/b4443a51336313a38b3caa627373488c5450a59df1aad26a3f197f46a16d17f1"
          },
          "type": "photo_camera",
          "title": "Foto pemilik toko",
          "description": "Ambil foto pemilik toko dengan jelas dan tidak terpotong",
          "placeholder": null,
          "validations": [
            {
              "type": "required",
              "error_message": "Harus diisi"
            }
          ]
        }
      ],
      "title": "Foto Pemilik Toko",
      "description": null
    },
    {
      "id": "task_with_store_owner_section",
      "info": {
        "icon": "https://content.payfazz.com/object/abbe92f80a439db37a1b81131feaa07c125b880af60b494cbffb02d79ecec807"
      },
      "type": "nextable_form",
      "child": [
        {
          "id": "owner_and_business_image",
          "info": {
            "hint": "",
            "icon": "https://content.payfazz.com/object/abbe92f80a439db37a1b81131feaa07c125b880af60b494cbffb02d79ecec807",
            "instruction_text": "",
            "instruction_image": "https://content.payfazz.com/object/f5ee4ca2cb8257245379bfd2cccdc642486a63ca47d2b5dfc5c14676322367dd"
          },
          "type": "photo_camera",
          "title": "Foto dengan pemilik toko",
          "description": "Ambil foto bersama pemilik toko di depan toko dengan jelas dan tidak terpotong",
          "placeholder": null,
          "validations": [
            {
              "type": "required",
              "error_message": "Harus diisi"
            }
          ]
        }
      ],
      "title": "Foto dengan Pemilik Toko",
      "description": null
    }
  ],
  "title": "Data Agen",
  "description": "Masukkan data agen yang akan diakuisisi"
}`
