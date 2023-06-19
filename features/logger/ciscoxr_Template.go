package logger

// ciscoxrTacacsTemplate is the template for CiscoXR AAA Tacacs.
const ciscoxrTemplate = `
!
service timestamps log datetime msec show-timezone
logging trap informational
logging archive
 device harddisk
 severity informational
 file-size 10
 archive-size 100
 archive-length 52
!
logging console disable
logging monitor informational
logging buffered 10000000
logging buffered informational
logging facility local1
{{ range .LoggerServerIPS }}
logging {{ . }} vrf default severity debugging port default
{{ end }}{{/* end range VIPs */}}
logging source-interface {{ .ManagementInterface }}
logging hostnameprefix {{ .Hostname }}
!
`
