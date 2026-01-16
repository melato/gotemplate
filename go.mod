module melato.org/gotemplate

go 1.23

replace (
	melato.org/command => ../command
	melato.org/yaml => ../yaml
)

require (
	melato.org/command v0.0.0-00010101000000-000000000000
	melato.org/yaml v0.0.0-00010101000000-000000000000
)

require gopkg.in/yaml.v2 v2.4.0 // indirect
