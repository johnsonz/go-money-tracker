$(function() {
    //datepicker setting
    $("#datepicker").datepicker({
        changeMonth: true,
        changeYear: true,
        dateFormat: "yy-mm-dd"
    });
    //get subcategory when select category
    $("#category").change(function() {
        $.getJSON("/getsubcategory?id=" + $(this).val(), function(data) {
            var options = '';
            $.each(data, function(key, val) {
                options += '<option value="' + val.ID + '" >' + val.Name + '</option>';
            });
            $("#subcategory").html(options);
        });
    });
    $(".smallimg").click(function(e) {
        $("body").find("#bigimg").remove();
        $("body").append('<p id="bigimg"><img src="' + this.src + '" alt="" /></p>');
        $(this).stop().fadeTo('slow', 0.5);
        $("#bigimg").fadeIn('fast');
        var w = document.documentElement.clientWidth;
        var h = document.documentElement.clientHeight;
        var top = (h - $("#bigimg").height()) / 2;
        var left = (w - $("#bigimg").width()) / 2;
        if (top < 0) {
            top = 0;
        }
        if (left < 0) {
            left = 0;
        }
        $("#bigimg").css({
            top: top,
            left: left
        });
        $("#bigimg").click(function() {
            $("body").find("#bigimg").remove();
        });
    });
    $("#navbar ul li").click(function() {
        Cookies.set("active", $(this).attr("name"));
    });
    setActiveNav();
});

function setActiveNav() {
    var active = Cookies.get("active");
    $("#navbar .active").removeClass("active");
    if (active == "" || active == undefined) {
        $("li[name='cate']").addClass("active");
    } else {
        $("li[name='" + active + "']").addClass("active");
    }
}
