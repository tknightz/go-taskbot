package trellobot

var LIST_MAP = map[string]string {
    "alltasks": "5f7addcd0ffbaa4f351c2226",
    "specifications": "5f7addcd0ffbaa4f351c2227",
    "todo": "5f7addcd0ffbaa4f351c2228",
    "doing": "5f87d86ba339027c8ae0697d",
    "edit": "5f87d86daad8aa44d0c3f8b2",
    "qc": "5f884b82afda200910ddd1ff",
    "done": "5f884b875c458316ac520816",
    "blocked": "5f884b8d3fc5fe57920c9965",
}

var LABELS_MAP = map[string]string {
    "1": "5f9e66cc4ea5040d5a8f6d42",
    "2": "5f9e66e4f988b3628674d45a",
    "3": "5f9e66ee29817880a2c60e4c",
}

var CUSTOMFIELDS_MAP = map[string]string {
    "mr": "5f87da9f9506437a7e9c6149",
}

var TYPELIST_MAP = map[string]string {
    "": "5f7addcd0ffbaa4f351c2228",
    "1": "5f7addcd0ffbaa4f351c2228",
    "2": "5f87d86daad8aa44d0c3f8b2",
    "3": "5f87d86ba339027c8ae0697d",
}


const (
    KEY = "key"
    TOKEN = "token"
    DESC = "desc"
    CALLBACK_URL = "https://test.com"
)
