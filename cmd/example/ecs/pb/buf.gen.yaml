version: v1

managed:
  enabled: true

plugins:
  - plugin: buf.build/protocolbuffers/go
    out: ./gen
    opt:
      - paths=source_relative

  - plugin: buf.build/community/mitchellh-go-json
    out: ./gen
    opt:
      - paths=source_relative

  - plugin: buf.build/community/chrusty-jsonschema:v1.4.1
    out: ./gen
    opt:
      - paths=source_relative
      - json_fieldnames
      - disallow_additional_properties
      - enforce_oneof
      - file_extension=jsonschema

  - plugin: buf.build/community/planetscale-vtprotobuf:v0.5.0
    out: ./gen
    opt:
      - paths=source_relative