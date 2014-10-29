package statik

import (
		"github.com/rakyll/statik/fs"
)

func init() {
	data := "PK\x03\x04\x14\x00\x08\x00\x00\x00\xc8)]E\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00	\x00\x00\x00test.html<!doctype html>\n<html>\n<head>\n  <title></title>\n  <link rel=\"stylesheet\" href=\"//cdn.jsdelivr.net/flat-ui/2.0/css/flat-ui.css\">\n  <style>\n    .inline{display:inline}\n    .filter{width:100%}\n    .container{ width: 800px; margin: 50px auto;}\n  </style>\n</head>\n<body>\n  <div class=\"container\">\n    <div class=\"row\">\n    </div>\n    <div class=\"row\">\n        <p>\n        <h3>Profiles</h3>\n        <button type=\"button\" class=\"cpu btn btn-primary\">CPU</button>\n        <button type=\"button\" class=\"heap btn btn-primary\">Heap</button>\n        <label class=\"checkbox\" for=\"cumsort\">\n            <input type=\"checkbox\" checked=\"checked\" id=\"cumsort\">\n            Cumulative sort\n        </label>\n        </p>\n        <input type=\"text\" class=\"filter\" placeholder=\"Filter by regex...\">\n        <div class=\"row\"><pre class=\"results\"></pre></div>\n      </div>\n  </div>\n  <script src=\"//ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min.js\"></script>\n  <script type=\"text/javascript\">\n    get(\"heap\");\n    $(\".cpu\").on(\"click\", function() {\n      get(\"profile\");\n    });\n    $(\".heap\").on(\"click\", function() {\n      get(\"heap\");\n    });\n    function get(name) {\n      $('.results').html('Loading, be patient...')\n      var f = $('.filter').val();\n      $.get('/p?profile=' + name + '&filter=' + f, function(data) {\n        $('.results').html(data);\n      });\n    };\n  </script>\n</body>\n</html>\nPK\x07\x08`a\x13Yk\x05\x00\x00k\x05\x00\x00PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xc8)]E`a\x13Yk\x05\x00\x00k\x05\x00\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x00\x00\x00\x00test.htmlPK\x05\x06\x00\x00\x00\x00\x01\x00\x01\x007\x00\x00\x00\xa2\x05\x00\x00\x00\x00"
	fs.Register(data)
}
