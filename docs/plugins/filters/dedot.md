### Fluentd Filter plugin to de-dot field name for elasticsearch.
#### More info at https://github.com/lunardial/fluent-plugin-dedot_filter

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| nested | bool | No | False | Will cause the plugin to recurse through nested structures (hashes and arrays), and remove dots in those key-names too. <br> |
| separator | string | No | _ | Separator <br> |
