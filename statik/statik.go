package statik

import (
		"github.com/rakyll/statik/fs"
)

func init() {
	data := "PK\x03\x04\x14\x00\x08\x00\x00\x00\xdd0]E\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00	\x00\x00\x00test.html<!doctype html>\n<html>\n<head>\n  <title></title>\n  <link rel=\"stylesheet\" href=\"//cdn.jsdelivr.net/flat-ui/2.0/css/flat-ui.css\">\n  <style>\n    .inline{display:inline}\n    .filter{width:100%}\n    .container{ width: 800px; margin: 50px auto;}\n  </style>\n</head>\n<body>\n  <div class=\"container\">\n    <div class=\"row\">\n    </div>\n    <div class=\"row\">\n        <h3>Profiles</h3>\n        <p>\n        <button type=\"button\" class=\"cpu btn btn-primary\">CPU</button>\n        <button type=\"button\" class=\"heap btn btn-primary\">Heap</button>\n        </p>\n        <input type=\"text\" class=\"filter\" placeholder=\"Focus by keyword...\">\n        <div class=\"row\"><pre class=\"results\"></pre></div>\n      </div>\n  </div>\n  <script src=\"//ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min.js\"></script>\n  <script type=\"text/javascript\">\n    var profile = \"heap\";\n    refresh();\n    $(\".cpu\").on(\"click\", function() {\n      profile = \"profile\";\n      refresh();\n    });\n    $(\".heap\").on(\"click\", function() {\n      profile = \"heap\";\n      refresh();\n    });\n    function refresh() {\n      $('.results').html('Loading, be patient... CPU profile takes 30 seconds.')\n      var f = $('.filter').val();\n      $.get('/p?profile=' + profile + '&filter=' + f, function(data) {\n        $('.results').html(data);\n      });\n    };\n  </script>\n</body>\n</html>\nPK\x07\x08\xa4\x92\x02\x902\x05\x00\x002\x05\x00\x00PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xdd0]E\xa4\x92\x02\x902\x05\x00\x002\x05\x00\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x00\x00\x00\x00test.htmlPK\x05\x06\x00\x00\x00\x00\x01\x00\x01\x007\x00\x00\x00i\x05\x00\x00\x00\x00"
	fs.Register(data)
}
