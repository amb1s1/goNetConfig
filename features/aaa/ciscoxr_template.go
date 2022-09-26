package aaa

// ciscoxrTacacsTemplate is the template for CiscoXR AAA Tacacs.
const ciscoxrTemplate = `
!
usergroup priv15
 taskgroup root-lr
 taskgroup cisco-support
!
username netops
 group root-lr
 group cisco-support
 secret 5 {{ .Secret }}
!
`
