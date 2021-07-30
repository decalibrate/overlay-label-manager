var $variables = jQuery("#variables");
var $templates = jQuery("#templates");
var $tmpl_variables = jQuery(".variables_row_template");
var $tmpl_templates = jQuery(".templates_row_template");

var variables = {};
var templates = {};
var configs = {};
var preferences = {};
var preference_defaults = {
    "enable_delete_warnings":true,
    "enable_codemirror":true
};

var codemirrorLodaing = false;
var codemirrorLoaded = false;

var templateCodeMirrors = [];


function addCodeMirrorToTemplateTextArea($ta) {
    if (codemirrorLoaded && checkPreference("enable_codemirror")) {
        var t = $ta.parent();

        if (!t.find("codemirror-index").length) {

            templateCodeMirrors.push(
                CodeMirror.fromTextArea($ta[0], {
                    lineNumbers: true,
                    lineWrapping: true,
                    mode: "liquid",
                    matchBrackets: true,
                    autoRefresh:true,
                    scrollbarStyle: "overlay"
                })
            );

            templateCodeMirrors[templateCodeMirrors.length -1].on("change", inputChangeCodeMirror, {target:$ta[0]});

            $ta.attr("codemirror-index", templateCodeMirrors.length -1);

            t.find(".CodeMirror-wrap")
                .addClass("ps-0")
                .addClass("pt-0")
                .addClass("pb-0")
                .addClass("form-control");
        }
    }
}

function onCodeMirrorLoaded(){

    if (checkPreference("enable_codemirror")) {
        jQuery("#templates textarea[data-template]").each(function() {
            var ta =  jQuery(this);
            var t = jQuery(this).parent();

            if (!ta.is("[codemirror-index]")) {
                addCodeMirrorToTemplateTextArea(ta);
            }

            var cm = ta.attr("codemirror-index");
            var cm_wrap = t.find(".CodeMirror-wrap")

            cm_wrap.prop("hidden", false);
            ta.prop("hidden", true)

            if (ta.hasClass("is-invalid")) {
                cm_wrap.addClass("is-invalid");
            } else if (ta.hasClass("is-valid")) {
                cm_wrap.addClass("is-valid");
            }

            ta
                .removeClass("is-invalid")
                .removeClass("is-valid");

            templateCodeMirrors[cm].setValue(ta.val());
            
        });
    }
};

function disableCodemirror() {
    if (codemirrorLoaded) {
        jQuery("#templates textarea[data-template][codemirror-index]").each(function() {
            var ta = jQuery(this);
            var cm_index = ta.attr("codemirror-index");
            var cm_wrap = ta.parent().find(".CodeMirror-wrap");

            ta.prop("hidden",false).show();
            cm_wrap.prop("hidden",true);

            ta.val(templateCodeMirrors[cm_index].getValue());

            if (cm_wrap.hasClass("is-invalid")) {
                ta.addClass("is-invalid");
            } else if (cm_wrap.hasClass("is-valid")) {
                ta.addClass("is-valid");
            }

            cm_wrap
                .removeClass("is-invalid")
                .removeClass("is-valid");
        });
    }
}

function enableCodemirror() {

    if (!codemirrorLodaing && !codemirrorLoaded) {
        var styles = [
            "https://cdn.jsdelivr.net/npm/codemirror@5.62.0/lib/codemirror.css",
            "https://cdn.jsdelivr.net/npm/codemirror@5.62.0/addon/scroll/simplescrollbars.css",
            "/codemirror-liquid-mode.css"
        ];

        var js = [
            "https://cdn.jsdelivr.net/npm/codemirror@5.62.0/lib/codemirror.js",
            "https://cdn.jsdelivr.net/npm/codemirror@5.62.0/addon/scroll/simplescrollbars.js",
            "https://cdn.jsdelivr.net/npm/codemirror@5.62.0/addon/display/autorefresh.js",
            "https://cdn.jsdelivr.net/npm/codemirror@5.62.0/mode/htmlmixed/htmlmixed.js",
            "/codemirror-liquid-mode.js"
        ];
        
        //    <link rel="stylesheet" href="/codemirror-liquid-mode.css">
        //    <script src="https://cdn.jsdelivr.net/npm/codemirror@5.62.0/lib/codemirror.js"></script>

        for (var i = 0; i < styles.length; i++) {
            var s =  document.createElement("link");
            s.setAttribute("rel","stylesheet");
            s.href = styles[i];
            document.head.appendChild(s);
        }

        var j = -1;

        var onload = function() {
            j++;
            if (j < js.length) {
                var s =  document.createElement("script");
                s.src = js[j];
                document.body.appendChild(s);

                s.onload = onload
            } else {
                codemirrorLoaded = true;
                codemirrorLodaing = false;

                onCodeMirrorLoaded();
            }
        }

        onload();
    } else {
        onCodeMirrorLoaded();
    }
}


function buildQS(data) {
    var qs = "";
    for (var p in data) {
        if (p != "name" && data.hasOwnProperty(p)) {
            var e = data[p];
            qs += p + "=" +  encodeURIComponent(e || "") + "&";
        }
    }
    qs = qs.slice(0,-1);

    return qs;
}

function numberElementFormat($p) {

    var num = $p.find('input[type="number"]');


    if (num.length) {
        var add = '<button class="btn btn-outline-secondary" type="button" data-bs-toggle="tooltip" data-add title="+1" aria-label="increment counter variable by 1"><i class="bi-plus"></i></button>';
        var sub = '<button class="btn btn-outline-secondary" type="button" data-bs-toggle="tooltip" data-sub title="-1" aria-label="decrment counter variable by 1"><i class="bi-dash"></i></button>';

        num.each(function() {
            var n = jQuery(this);
            var h = n.clone();
            n.replaceWith(jQuery('<div class="input-group input-group-sm">' + add + sub + '</div>').prepend(h));

        });

    }
}

function setVariableData(v, data) {
    v.find("[data-original_name]").val(data.name || "");
    v.find("[data-type]").val(data.type || "");
    v.find("[data-name]").val(data.name || "");
    v.find("[data-value]").val(data.value || "0");
    v.find("[data-goal]").val(data.goal || "");
    v.find("[data-completion_text]").val(data.completion_text || "");
}

function activateTooltips($el) {
    $el.find('[data-bs-toggle="tooltip"]').each(function(k,v) {
        var tt = new bootstrap.Tooltip(v, {delay:{show:1000}, trigger:"hover manual"});
        $el.on("click", function() {
            tt.hide();
        }).on("blur", function() {
            tt.hide();
        });
    });
}

function setVariableType(v, t, f) {
    v.find("[fields]").removeAttr("hidden");
    v.find("[data-type]").val(t);
    var icon = v.find("[data-ts] [data-to=\"" + t + "\"] i");
    var vi = v.find("[variable-icon]");
    vi
        .removeClass("bi-dash")
        .addClass(icon.attr("class"))
        .parent()
        .removeAttr("hidden");

    title_text = icon.parent().text().trim();

    vi.attr("data-bs-original-title", title_text + " Variable");

    var tt = bootstrap.Tooltip.getInstance(vi[0]);
    if(tt) {
        tt.update();
    }

    v.find("[data-ts]").remove();
    v.find("[data-s], [data-c], [show-more]").parent().removeClass("d-none");

    v.find("[data-used-by]").each(function() {
        el = jQuery(this)
        if ((el.attr("data-used-by") || "").match(t) || (el.attr("data-used-by") || "").indexOf("*") > -1) {

        } else {
            el.remove();
        }
    });

    var more = v.find("[show-more]");
    var more_content = v.find(".collapse");

    if (more_content.children().length == 0) {
        more.remove();
    }

    if (f) {
        v.find("[data-name]")[0].focus();
    }
}

function addNewVariable(data) {
    // create html components
    var v = $tmpl_variables.clone().removeAttr("hidden");

    numberElementFormat(v);

    // force auto complete off
    v.find("input").prop("autocomplete", "off");


    if (data) {
        setVariableData(v,data);
        setVariableType(v,data.type);
    } else {
        v.find(".dropdown-item").on("click", function(ee) {
            e = jQuery(this);
            setVariableType(v,e.attr("data-to"), true);
        });
    }

    activateTooltips(v);

    var plusMinus = v.find("button[data-add], button[data-sub]")

    var more = v.find("[show-more]");
    var chevron = v.find("i.bi-chevron-down");
    var more_content = v.find(".collapse");

    var bs_collapse = new bootstrap.Collapse(more_content, {toggle:false});

    more.click(function(e) {
        e.preventDefault();

        if(more_content.hasClass("show")) {
            bs_collapse.hide();
            chevron
                .addClass("bi-chevron-down")
                .removeClass("bi-chevron-up");
        } else {
            bs_collapse.show();
            chevron
                .addClass("bi-chevron-up")
                .removeClass("bi-chevron-down");
        }
    });

    $variables.append(v);
    
    
    if (typeof data == "undefined" || !data) {
        addElementExpand(v,v.find("[data-name]:eq(0)"));
    }
}

function deleteVariable($t) {
    // remove
    var v = $t.parents(".variables_row_template");

    var data = {};

    data.name = v.find("[data-original_name]").val() || "";

    if (data.name && variables[data.name]) {

        var req = jQuery.ajax({
            url:"/v/" + data.name,
            method:"DELETE"
        });

        delete variables[data.name];
    }

    deleteElementCollapse(v);
}

function cancelVariableChanges($t) {
    var v = $t.parents(".variables_row_template");
    
    var originalName = v.find("[data-original_name]").val() || "";
    var data = variables[originalName];

    if (data) {
        setVariableData(v,data);

        v.find("input, textarea")
            .removeClass("is-valid")
            .removeClass("is-invalid");
        
        v.find("[data-s], [data-c]")
            .prop("disabled",true)
            .addClass("btn-outline-primary")
            .removeClass("btn-primary");
    }
}

function saveVariable($t) {
    var v = $t.parents(".variables_row_template");
    
    var data = {};

    // data.name = v.find("[data-name]").val() || "";
    data.type = v.find("[data-type]").val() || "";
    var originalName = v.find("[data-original_name]").val() || "";

    v.find("[data-used-by='" + data.type + "'] input, [data-used-by='*'] input").each(function() {
        var dv = jQuery(this);

        var k = Object.values(dv[0].attributes).filter(function(k,v){return k.nodeName.indexOf("data-") == 0});

        if (k.length > 0) {
            var key = k[0].nodeName.split("data-")[1];
            data[key] = dv.val() || "";
        }
    });

    if (data.name && data.type) {

        var qs = buildQS(data);

        var req = jQuery.ajax({
            url:"/v/" + data.name + "?" + qs,
            method:"PUT"
        }).done(function(e) {
            v.find("[data-original_name]").val(e.name);
            variables[e.name] = e;
        });

        if ( originalName != "" && originalName != data.name ) {
            var delreq = jQuery.ajax({
                url:"/v/" + originalName,
                method:"DELETE"
            }).done(function(e) {
                delete variables[originalName];
            });
        }


        // req.done(function(e) {
        //     console.log(e);
        // })
        // req.fail(function(e) {
        //     console.log(e);
        // })

        v.find("input, textarea")
        .removeClass("is-valid")
        .removeClass("is-invalid");
    
        v.find("[data-s], [data-c]")
            .prop("disabled",true)
            .addClass("btn-outline-primary")
            .removeClass("btn-primary");
    }
}





function setTemplateData(t,data) {
    t.find("input[data-original_name]").val(data.name || "");
    t.find("input[data-name]").val(data.name || "");

    var ta = t.find("textarea[data-template]")

    if(checkPreference("enable_codemirror") && codemirrorLoaded) {
        templateCodeMirrors[ta.attr("codemirror-index")].setValue(data.template);
    } else {
        ta.val(data.template || "");
    }

    t.find("button[data-hidden]").prop("disabled",false);
    setHideShowButton(t.find("button[data-hidden] i"), data.hidden || false);
}

function addNewTemplate(data) {
    // create html components
    var t = $tmpl_templates.clone().removeAttr("hidden");

    // force auto complete off
    t.find("input").prop("autocomplete", "off");

    if (typeof data != "undefined") {
        setTemplateData(t,data);
    }

    activateTooltips(t);

    $templates.append(t);

    if (typeof data == "undefined") {
        addElementExpand(t,t.find("[data-name]:eq(0)"));
    }

    addCodeMirrorToTemplateTextArea(t.find("textarea"));
}

function deleteTemplate($t) {
    // remove
    var t = $t.parents(".templates_row_template");

    var data = {};

    data.name = t.find("input[data-original_name]").val() || "";

    if (data.name) {

        var qs = buildQS(data);

        var req = jQuery.ajax({
            url:"/t/" + data.name + "?" + qs,
            method:"DELETE"
        });

        delete templates[data.name];
    }

    deleteElementCollapse(t);
}

function cancelTemplateChanges($t) {
    var t = $t.parents(".templates_row_template");

    var data = {};

    var originalName = t.find("[data-original_name]").val() || "";
    var data = templates[originalName];

    if (data) {
        setTemplateData(t,data);

        t.find("input, .CodeMirror.is-valid, .CodeMirror.is-invalid")
            .removeClass("is-valid")
            .removeClass("is-invalid")

        t.find("[data-s], [data-c]")
            .prop("disabled",true)
            .addClass("btn-outline-primary")
            .removeClass("btn-primary");
    }
}

function saveTemplate($t) {
    var t = $t.parents(".templates_row_template");

    var data = {};

    data.name = t.find("input[data-name]").val() || "";
    var originalName = t.find("[data-original_name]").val() || "";
    
    if (data.name) {

        var ta = t.find("textarea[data-template]");
        if (checkPreference("enable_codemirror")) {
            data.template = templateCodeMirrors[ta.attr("codemirror-index")].getValue() || "";
        } else {
            data.template = ta.val() || "";
        }

        if ( t.find("button[data-hidden] i.bi-eye-slash").length ) {
            data.hidden = true;
        }

        var req = jQuery.ajax({
            url:"/t/" + encodeURIComponent(data.name) + (data.hidden ? "?hidden=true" : ""),
            method:"POST",
            data: data.template,
            contentType: "text/plain"
        }).done(function(e) {
            templates[e.name] = e;
            setTemplateData(t,e);
        });
        

        if ( originalName != "" && originalName != data.name) {
            var delreq = jQuery.ajax({
                url:"/t/" + originalName,
                method:"DELETE"
            }).done(function(e) {
                delete templates[originalName];
            });
        }

        t.find("input, textarea, .CodeMirror.is-valid, .CodeMirror.is-invalid")
            .removeClass("is-valid")
            .removeClass("is-invalid");
    
        t.find("[data-s], [data-c]")
            .prop("disabled",true)
            .addClass("btn-outline-primary")
            .removeClass("btn-primary");
    }
}

function setHideShowButton($i, hide) {
    var $b = $i.parent("button:eq(0)");

    if (!hide) {
        // eye open
        $i.addClass('bi-eye-fill')
            .removeClass('bi-eye-slash');
        
        $b.addClass("btn-secondary")
            .removeClass("btn-outline-secondary");

        $b.parent().prop("title", "Showing content in label");
    } else {
        // eye slashed
        $i.removeClass('bi-eye-fill')
            .addClass('bi-eye-slash');

        $b.removeClass("btn-secondary")
            .addClass("btn-outline-secondary");
        
        $b.parent().prop("title", "Hiding content in label");
    }
}

function toggleTemplateHide($t) {
    var t = $t.parents(".templates_row_template");

    var originalName = t.find("[data-original_name]").val() || "";
    var data = templates[originalName];

    
    if (data.name) {

        var $i = t.find("button[data-hidden] i");
        var $b = $i.parent();

        if (!data.hidden) {
            // hide content on clicking button
            data.hidden = true;
        } else if (data.hidden) {
            // unhide content on clicking button
            delete data.hidden;
        }

        setHideShowButton($i, data.hidden || false);

        var req = jQuery.ajax({
            url:"/t/" + encodeURIComponent(data.name) + (data.hidden ? "?hidden=true" : ""),
            method:"PUT"
        }).done(function(e) {
            templates[e.name] = e;
        });
    }
}


function fillConfigs(data) {
    for (var c in data) {
        if (Object.hasOwnProperty.call(data, c)) {
            $c = jQuery("#configs [data-" + c + "]").val(data[c]);
            configs[c] = data[c];
        }
    }

    jQuery("#version").text(configs.version).addClass("d-inline-block").prop("hidden",false);

    if(configs.updateAvailable) {
        jQuery(".update-prompt").addClass("d-inline-block").prop("hidden",false);
    }
}


function saveConfigurationChanges() {
    var c = jQuery("#configs");

    var data = {};

    data.port = parseInt(c.find("[data-port]").val() || configs.port, 10);
    data.labelDirectory = c.find("[data-labelDirectory]").val() || configs.labelDirectory;
    
    if (Object.keys(data).length > 0) {

        var req = jQuery.ajax({
            url:"/conf",
            method:"POST",
            data: JSON.stringify(data),
            contentType: "application/json"
        }).done(function(data) {

            if (data.port && configs.port != data.port) {
                // redirect
                window.location.hash = '#configuration';
                window.location.port = data.port;
            } else {
                fillConfigs(data);

                c.find("input, textarea")
                    .removeClass("is-valid")
                    .removeClass("is-invalid");
            
                c.find("[data-s], [data-c]")
                    .prop("disabled",true)
                    .addClass("btn-outline-primary")
                    .removeClass("btn-primary");
            }
        });

    }
}
function cancelConfigurationChanges($t) {
    var c = jQuery("#configs");


    fillConfigs(configs);

    c.find("input, textarea")
        .removeClass("is-valid")
        .removeClass("is-invalid");
    
    c.find("[data-s], [data-c]")
        .prop("disabled",true)
        .addClass("btn-outline-primary")
        .removeClass("btn-primary");
}



function validateVTName($t) {
    if ($t.val().match("^[a-zA-Z_][0-9a-zA-Z_]{0,19}$")) {
        var $vs = $t.parents("#templates, #variables").find(".vt-line:not('[hidden]')").find("input[type=text][data-name]").not($t);
        if ( $vs.filter(function(){ return jQuery(this).val() == $t.val() }).length == 0 ) {
            return true
        } else {
            // console.log("another variable with this name");
        }
    } else {
        // console.log("formatting issue");
    }

    return false
}

function validateNumber(el) {
    if (el.val().match("^(-?[0-9]+(\.[0-9]+)?)?$")) {
        return true
    }

    return false
}

function validatePort(el) {
    if (el.val().match("^(-?[0-9]+(\.[0-9]+)?)?$")) {
        if(el.val() > 0 && el.val() < 65536 /* 2 ^ 16 */) {
            return true
        }
    }

    return false

}

function setActiveTab() {
    var hash = document.location.hash + "";

    if (hash.indexOf("#documentation-") == 0) {
        var trigger = jQuery('a[href="' + document.location.hash + '"]');
        highlightDocumentationSelect({target:trigger[0]}, true);
        hash = "#documentation";
    }

    if (jQuery('#navbarNav a[href="' + hash + '"]').length) {
        var t = new bootstrap.Tab(jQuery('#navbarNav a[href="' + hash + '"]')[0]);
        t.show();
    }
}


function highlightDocumentationSelect(trigger, scrollTo) {

    var target = jQuery(trigger.target);

    if (!target.is("a[href^='#documentation-']")){
        target = target.parents("a");
    }

    var doc_header = jQuery(target.attr("href"));
    doc_header.find("h4, h5, h6, p.fst-italic").eq(0).addClass("text-primary");
    window.setTimeout(function() { doc_header.find("h4, h5, h6, p.fst-italic").eq(0).removeClass("text-primary"); }, 1500);

    if (typeof trigger.preventDefault != "undefined") {

        window.setTimeout(function() {
            doc_header[0].scrollIntoView()
        },100);
        trigger.preventDefault();
    }
}


function addElementExpand($e, focusEl) {
    if(!$e.hasClass("collapse")) {
        $e.addClass("collapse");

        var bs_collapse = new bootstrap.Collapse($e[0], {toggle:false});

        $e[0].addEventListener('shown.bs.collapse', function (e) {
            if(e.target == $e[0]) {
                $e
                    .removeClass("collapse")
                    .removeClass("show");

                focusEl.focus();
                bs_collapse.dispose();
            }
        });

        window.setTimeout(function() {
            bs_collapse.show();
        }, 100);
    }
}


function deleteElementCollapse($e) {
    if(!$e.hasClass("collapse")) {
        $e.addClass("show").addClass("collapse");

        var bs_collapse = new bootstrap.Collapse($e[0], {toggle:false});

        $e[0].addEventListener('hidden.bs.collapse', function (e) {
            if(e.target == $e[0]) {
                bs_collapse.dispose();
                if ( $e.find(".CodeMirror").length ) {
                    var cm_index = $e.find("codemirror-index");
                    templateCodeMirrors[cm_index] = null;
                }
                $e.remove();
            }
        });

        window.setTimeout(function() {
            bs_collapse.hide();
        }, 200);
    }
}

function inputChangeCodeMirror(cm) {
    inputChange({target:cm.getTextArea()}, cm);
}

function inputChange(e, cm) {

    var $p, typ;
    var el = jQuery(e.target);
    
    if (el.is("textarea[codemirror-index]")) {
        if (checkPreference("enable_codemirror")) {
            if(!cm) {
                return
            }
            el = el.parent().find(".CodeMirror-wrap")
        }
    }

    el.removeClass("is-invalid");
    el.removeClass("is-valid");


    if (el.parents("#templates").length > 0) {
        $p = el.parents(".templates_row_template");
        typ = templates;
    } else if (el.parents("#configs").length > 0) {
        $p = el.parents("#configs");
        typ = configs;
    } else {
        $p = el.parents(".variables_row_template");
        typ = variables;
    }


    var isValid = true;
    if(el.is("[data-name]")) {
        isValid = validateVTName(el);
    } else if(el.is("[data-port]")) {
        isValid = validatePort(el);
    } else if (el.is("[type=number]")) {
        isValid = validateNumber(el)
    }

    var on, obj;

    if (typ === configs) {
        obj = configs;
    } else {
        on = $p.find("[data-original_name]").val();
        obj = typ[on];
    }
    
    var val;
    if (cm) {
        val = cm.getValue();
    } else {
        val = el.val();
    }

    if (isValid) {


        var prop = Object.values(jQuery(e.target)[0].attributes).filter(function(v) {
            return v.nodeName.indexOf("data-") == 0
        })[0].nodeName.split("data-")[1];


        if (obj) {
            objpropdefined = obj.hasOwnProperty(prop);
            if ((!objpropdefined && val) || objpropdefined && obj[prop] != val) {
                if (prop == "n" && val != on && typ[val]) {
                    el.addClass("is-invalid");
                } else {
                    el.addClass("is-valid");
                }
            }
        } else if (prop == "n" && val != on && typ[val]) {
            el.addClass("is-invalid");
        } else {
            el.addClass("is-valid");
        }
    } else {
        el.addClass("is-invalid");
    }

    var $c = $p.find("button[data-c]");
    var $s = $p.find("button[data-s]");
    
    if ($p.find(".is-invalid").length > 0)  {
        if (obj) {
            $c.prop("disabled", false)
                .removeClass("btn-outline-primary")
                .addClass("btn-primary");
        }

        $s.prop("disabled", true)
            .addClass("btn-outline-primary")
            .removeClass("btn-primary");
    } else if ($p.find(".is-valid").length > 0) {
        if (obj) {
            $c.prop("disabled", false)
                .removeClass("btn-outline-primary")
                .addClass("btn-primary");
        }

        $s.prop("disabled", false)
            .removeClass("btn-outline-primary")
            .addClass("btn-primary");
    } else {
        $c.prop("disabled", true)
            .addClass("btn-outline-primary")
            .removeClass("btn-primary");

        $s.prop("disabled", true)
            .addClass("btn-outline-primary")
            .removeClass("btn-primary");
    }
}

jQuery(document).on("input","#variables input, #templates input, #templates textarea[data-template], #configs input", inputChange);

var modal_delete_target = null;
var deleteModal = jQuery("#delete_modal");
var deleteModal_bs = new bootstrap.Modal(deleteModal[0], {});

deleteModal[0].addEventListener('hide.bs.modal', function (event) {
    modal_delete_target = null;
});


function modalDeletePrep($el) {
    modal_delete_target = $el;
    
    var p,on ;
    deleteModal.find("[del-label-hint], [del-variable-hint]").prop("hidden",true);


    if ( $el.parents("#templates").length > 0 ) {
        p = $el.parents(".templates_row_template");
        on = p.find("[data-original_name]").val();
        deleteModal.find("[del-type]").text("label");
        if(on) {
            deleteModal.find("[del-label-hint]").prop("hidden",false);
        }
        deleteModal.find("[del-var-type]").html('<i class="bi-tags"></i>');
    } else {
        p = $el.parents(".variables_row_template");
        on = p.find("[data-original_name]").val();
        deleteModal.find("[del-type]").text("variable");
        if(on) {
            deleteModal.find("[del-variable-hint]").prop("hidden",false);
        }
        deleteModal.find("[del-var-type]").html(p.find("[variable-icon]").clone()); 
    }

    var n = on || p.find("[data-name]").val() || "unnamed";
    deleteModal.find("[del-name]").text(n);
}

function modalDeleteAction($el) {
    if ( modal_delete_target != null) {
        $el = modal_delete_target;
    }

    if ($el) {
        if ( $el.parents("#templates").length > 0 ) {
            deleteTemplate($el);
        } else {
            deleteVariable($el);
        }
    }
    modal_delete_target = null;

    deleteModal_bs.hide();
}

jQuery(document).on("click","button[data-d]", function(e) {
    var $el = jQuery(this);

    var p_row = $el.parents(".templates_row_template, .variables_row_template");

    var isWanted = (p_row.find("[data-original_name]").val() || p_row.find(".is-valid").length);

    if (!isWanted || !checkPreference("enable_delete_warnings") ) {
        modalDeleteAction($el)
    } else {
        modalDeletePrep($el);
        deleteModal_bs.show();
    }
});


jQuery(document).on("click","button[data-c]", function(e) {
    var $el = jQuery(this);
    if ( $el.parents("#templates").length > 0 ) {
        cancelTemplateChanges($el);
    } else if ( $el.parents("#preferences").length > 0 ) {
        cancelPreferenceChanges();
    } else if ( $el.parents("#configs").length > 0 ) {
        cancelConfigurationChanges($el);
    } else {
        cancelVariableChanges($el);
    }
});

jQuery(document).on("click","button[data-s]", function(e) {
    var $el = jQuery(this);
    if ( $el.parents("#templates").length > 0 ) {
        saveTemplate($el);
    } else if ( $el.parents("#preferences").length > 0 ) {
        savePreferenceChanges();
    } else if ( $el.parents("#configs").length > 0 ) {
        saveConfigurationChanges($el);
    } else {
        saveVariable($el);
    }
});

jQuery(document).on("click","button[data-hidden]", function(e) {
    var $el = jQuery(this);
    if ( $el.parents("#templates").length > 0 ) {
        toggleTemplateHide($el);
    }
});


function checkPreference(p){
    if (preferences.hasOwnProperty(p)) {
        return preferences[p];
    } else {
        return preference_defaults[p];
    }
}

function _jsonp(t) {
    var endpoint;

    switch (t) {
        case "variables":
            jQuery(".variables-container").prop("hidden", true);
            jQuery(".variables-spinner").prop("hidden",false);
            endpoint = "v";
            break;
        case "templates":
            jQuery(".templates-container").prop("hidden", true);
            jQuery(".templates-spinner").prop("hidden",false);
            endpoint = "t";
            break;
        case "configs":
            jQuery(".configs-container").prop("hidden", true);
            jQuery(".configs-spinner").prop("hidden",false);
            endpoint = "conf";
            break;
    }

    if (endpoint) {
        var scr = document.createElement("script");
        var a = new Date();
        scr.src = '/' + endpoint + '?jsonp=1&cbn=' + t + '&cbp=window._jsonp&ts=' + a.getTime();
        document.body.appendChild(scr)
    }
};

_jsonp.variables = function(data) {
    for (var v in data) {
        if (Object.hasOwnProperty.call(data, v)) {

            variables[data[v].name] = data[v];
            addNewVariable(data[v]);
        }
    }
    jQuery(".variables-container").prop("hidden",false);
    jQuery(".variables-spinner").prop("hidden", true);

    _jsonp.callbackTimeouts.variables = null;
};

_jsonp.templates = function(data) {
    for (var t in data) {
        if (Object.hasOwnProperty.call(data, t)) {

            templates[data[t].name] = data[t];
            addNewTemplate(data[t]);
        }
    }
    jQuery(".templates-container").prop("hidden",false);
    jQuery(".templates-spinner").prop("hidden", true);

    _jsonp.callbackTimeouts.templates = null;
};

_jsonp.configs = function(data) {
    fillConfigs(data);
    jQuery(".configs-container").prop("hidden",false);
    jQuery(".configs-spinner").prop("hidden", true);

    _jsonp.callbackTimeouts.configs = null;
};

_jsonp.callbackTimeouts = {
    configs: 1,
    templates: 1,
    variables: 1
};

function setPreferences(el) {
    if (window.localStorage && window.localStorage.getItem) {
        var prefs = JSON.parse(window.localStorage.getItem("__olm_preferences") || "{}");
        var $pr = jQuery("#preferences");
        for (var p in preference_defaults) {
            if (Object.hasOwnProperty.call(preference_defaults, p)) {

                var $p = $pr.find("#" + p);

                if (typeof el != "undefined") {
                    delete preferences[p];
                    if ($p.attr("type") == "checkbox") {
                        if($p.is(":checked")) {
                            if ( preference_defaults[p] !== true ) {
                                preferences[p] = true;
                            }
                        } else {
                            if ( preference_defaults[p] !== false ) {
                                preferences[p] = false;
                            }
                        }
                    }
                } else if (Object.hasOwnProperty.call(prefs, p)) {
                    preferences[p] = prefs[p];
                    if ($p.attr("type") == "checkbox") {
                        $p.prop("checked", prefs[p]);
                    }
                }
            }
        }

        if (Object.keys(preferences).length){
            window.localStorage.setItem("__olm_preferences",JSON.stringify(preferences));
        } else {
            window.localStorage.removeItem("__olm_preferences");
        }

        var changed = jQuery(this);

        if (changed.is("input#enable_codemirror")) {
            if(changed.is(":checked")) {
                enableCodemirror();
            } else {
                disableCodemirror();
            }
        }
    }
}

jQuery(document).on("click", "button[data-add], button[data-sub]", function(e) {
    var b = jQuery(this);
    var i = b.parent().find("input");

    if (i.val() == "") {
        i.val("0");
    }
    
    if ( b.is("[data-add]") || b.parents("[data-add]").length) {
        i.val(parseFloat(i.val()) + 1);
        i.trigger("input");
    } else if ( b.is("[data-sub]") || b.parents("[data-sub]").length) {
        i.val(parseFloat(i.val()) - 1);
        i.trigger("input");
    }
});




jQuery(document).ready(function() {


    setActiveTab();
    jQuery("#documentation a[href^='#documentation-']").on("click", highlightDocumentationSelect);

    if (window.localStorage && window.localStorage.getItem) {
        setPreferences();
        
        jQuery("#preferences input").on("change", setPreferences);
    } else {
        jQuery("#preferences input").prop("disabled", true);
        jQuery("#localStorage-fail").prop("hidden",false)
    }

    window.setTimeout(function(){
        for (var to in _jsonp.callbackTimeouts) {
            if (_jsonp.callbackTimeouts[to] != null) {
                jQuery("." + to + "-spinner").prop("hidden", true);
                jQuery("." + to + "-spinner-fail").prop("hidden", false);
            }
        }
    }, 3000);

    _jsonp("variables");
    _jsonp("templates");
    _jsonp("configs");

    activateTooltips(jQuery("#configuration"));
    activateTooltips(jQuery(".footer"));

    if (checkPreference("enable_codemirror")) {
        enableCodemirror();
    }
    
})


// var input = document.getElementById("input");
// var output = document.getElementById("output");
// var socket = new WebSocket("ws://" + window.location.hostname + (window.location.port? ":" + window.location.port : "") + "/ws");


// socket.onopen = function() {
//     output.innerHTML += "Status: Connected\n";
// };

// socket.onmessage = function(e) {
//     output.innerHTML += "Server: " + e.data + "\n";
// };

// function send() {
//     socket.send(input.value);
//     input.value = "";
// }