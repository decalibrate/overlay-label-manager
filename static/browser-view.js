function liquifyLoad() {
    var qs = window.location.search;
    var name = (window.location.pathname + "").split("/bv/")[1].split('/')[0];

    var params = {};

    if (qs[0] == "?") {
        qs = qs.substring(1);
    }

    qs = qs.split("&");

    function repeatedDecode(s) {
        var n = 0;
        while (decodeURIComponent(s) != s && n < 10) {
            try {
                ksv = decodeURIComponent(s);
                n++;
            } catch (e) {
                n = 10;
            }
        }
        return s;
    }

    for (var i = 0; i < qs.length; i++) {
        var kv = qs[i].split("=");
        kv[0] = repeatedDecode(kv[0]);
        if (kv.length > 1) {
            kv[1] = repeatedDecode(kv[1]);
        } else {
            kv[1] = "";
        }
        if (kv[0]) {
            params[kv[0]] = kv[1];
        }
    }

    if (name) {
        params.name = name;
    }

    console.log("Query String: " + qs);
    console.log("Params: " + JSON.stringify(params));


    var template = {};
    var variables = {};
    var label = "";


    function processTemplate() {

        if (template.hidden) {
            return "";
        }


        var Liquid = window.liquidjs.Liquid
        var engine = new Liquid({
            extname: '.html',
            cache: true
        });

        engine
            .parseAndRender(template.template, variables)
            .then(function(html) {
                document.getElementById("label").innerHTML = html
            });
    }

    if (params.name) {
        var request = new XMLHttpRequest();
        request.open('GET', "/t/" + encodeURIComponent(params.name) + "?bv=true", true);

        request.onload = function() {
            if (this.status >= 200 && this.status < 400) {
                // Success!
                var data = JSON.parse(this.response);

                template = data.t || {};
                console.log("Template: " + JSON.stringify(template));

                for (var vk = 0; vk < (data.v || []).length; vk++) {
                    var v = data.v[vk];
                    variables[v.name] = v;
                }

                console.log("Variables: " + JSON.stringify(variables));

                processTemplate();
            } else {
                // We reached our target server, but it returned an error
                console.log("could not get template:", data);

            }
        };

        request.onerror = function() {
            // There was a connection error of some sort
            console.log("could not get template:", data);
        };

        request.send();
    }

}

var liquify = document.createElement("script");
liquify.src = "https://cdn.jsdelivr.net/npm/liquidjs/dist/liquid.browser.min.js";
liquify.onload = liquifyLoad;
document.body.appendChild(liquify);