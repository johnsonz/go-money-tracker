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
    $('#Modal').on('show.bs.modal', function(event) {
        var element = $(event.relatedTarget) // element that triggered the modal
        var ep=element.parent().parent();
        var name=ep.find("span[name='catename']").html();
        var time=ep.find("span[name='catectime']").html();
        var by=ep.find("span[name='catecby']").html();
        var modal = $(this)
        modal.find('.modal-body #catename').val(name);
        modal.find('.modal-body #createdtime').val(time);
        modal.find('.modal-body #createdby').val(by);
    })
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
