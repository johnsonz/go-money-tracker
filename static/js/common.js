$(function() {
    //datepicker setting
    $(".datepicker").datepicker({
        changeMonth: true,
        changeYear: true,
        dateFormat: "yy-mm-dd"
    });
    //get subcategory when select category
    $("#category,#updatedcategory").change(function() {
        var options = '';
        var flag = $(this).attr("id");

        $.getJSON("/getsubcategory?id=" + $(this).val(), function(data) {
            $.each(data, function(key, val) {
                options += '<option value="' + val.ID + '" >' + val.Name + '</option>';
            });
            if (flag == "category") {
                $("#subcategory").html(options);
            } else {
                $("#updatedsubcategory").html(options);
            }
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
        var ep = element.parent().parent();
        var id = ep.find("span[name='cateid']").html();
        var name = ep.find("span[name='catename']").html();
        var time = ep.find("span[name='catectime']").html();
        var by = ep.find("span[name='catecby']").html();
        var modal = $(this)
        $('#updatedid').val(id);
        modal.find('.modal-body #catename').val(name);
        modal.find('.modal-body #createdtime').val(time);
        modal.find('.modal-body #createdby').val(by);
    });
    $('#Modal-Subcate').on('show.bs.modal', function(event) {
        var element = $(event.relatedTarget) // element that triggered the modal
        var ep = element.parent().parent();
        var id = ep.find("span[name='subcateid']").html();
        var name = ep.find("span[name='subcatename']").html();
        var time = ep.find("span[name='subcatectime']").html();
        var by = ep.find("span[name='subcatecby']").html();
        var modal = $(this)
        $('#updatedid').val(id);
        modal.find('.modal-body #subcatename').val(name);
        modal.find('.modal-body #createdtime').val(time);
        modal.find('.modal-body #createdby').val(by);
    });
    $('#Modal-Item').on('show.bs.modal', function(event) {
        var element = $(event.relatedTarget) // element that triggered the modal
        var ep = element.parent().parent();

        var id = ep.find("span[name='itemid']").html();
        var cate = ep.find("span[name='itemcate']").html();
        var subcate = ep.find("span[name='itemsubcate']").html();
        var store = ep.find("span[name='itemstore']").html();
        var addr = ep.find("span[name='itemaddr']").html();
        var pur = ep.find("span[name='itempurdate']").html();
        var amount = ep.find("span[name='itemamount']").html();
        var receipt = ep.find("span[name='itemreceipt']").html();
        var remark = ep.find("span[name='itemremark']").html();
        var time = ep.find("span[name='itemctime']").html();
        var by = ep.find("span[name='itemcby']").html();
        var modal = $(this)
        $('#updatedid').val(id);
        modal.find('.modal-body #updatedcategory').find("option[text='" + cate + "']").attr("selected", true);
        var cateid = modal.find('.modal-body #updatedcategory').val();
        $.getJSON("/getsubcategory?id=" + cateid, function(data) {
            var options = '';
            $.each(data, function(key, val) {
                options += '<option value="' + val.ID + '" >' + val.Name + '</option>';
            });
            $("#updatedsubcategory").html(options);
            modal.find('.modal-body #updatedsubcategory').find("option[text='" + subcate + "']").attr("selected", true);
        });

        modal.find('.modal-body #updatedstore').val(store);
        modal.find('.modal-body #updatedaddress').val(addr);
        modal.find('.modal-body #updatedpurchaseddate').val(pur);
        //modal.find('.modal-body #purchaseddatereceipt').val(amount);
        if (receipt != "None") {
            var remove = '<input type="button" value="Remove" class="btn btn-default btn-xs" id="rmRept"/>';
            modal.find('.modal-body #wrapreceip').html(receipt + remove);
        } else {
            modal.find('.modal-body #wrapreceip').html('');
        }
        $("#rmRept").click(function() {
            $.post("/rmrept", {
                    id: id
                },
                function(data, status) {
                    if (data) {
                        modal.find('.modal-body #wrapreceip').html('');
                        ep.find("span[name='itemreceipt']").html('None');
                    }
                });
        });
        modal.find('.modal-body #purchaseddateremark').val(remark);
        modal.find('.modal-body #createdtime').val(time);
        modal.find('.modal-body #createdby').val(by);
    });
    $('#Modal-Detail').on('show.bs.modal', function(event) {
        var element = $(event.relatedTarget) // element that triggered the modal
        var ep = element.parent().parent();

        var id = ep.find("span[name='detailid']").html();
        var name = ep.find("span[name='detailname']").html();
        var price = ep.find("span[name='detailprice']").html();
        var quan = ep.find("span[name='detailquan']").html();
        var amount = ep.find("span[name='detailamount']").html();
        var lone = ep.find("span[name='detaillone']").html();
        var ltwo = ep.find("span[name='detailltwo']").html();
        var remark = ep.find("span[name='detailremark']").html();
        var time = ep.find("span[name='detailctime']").html();
        var by = ep.find("span[name='detailcby']").html();

        var modal = $(this)
        $('#updatedid').val(id);
        modal.find('.modal-body #updatedname').val(name);
        modal.find('.modal-body #updatedprice').val(price);
        modal.find('.modal-body #updatedquantity').val(quan);
        modal.find('.modal-body #updatedamount').val(amount);
        // modal.find('.modal-body #updatedlone').val(lone);
        // modal.find('.modal-body #updatedlone').val(ltwo);
        modal.find('.modal-body #updatedremark').val(remark);
        modal.find('.modal-body #createdtime').val(time);
        modal.find('.modal-body #createdby').val(by);
        if (lone != "None") {
            var remove = '<input type="button" value="Remove" class="btn btn-default btn-xs" id="rmlOne"/>';
            modal.find('.modal-body #wraplone').html(lone + remove);
        } else {
            modal.find('.modal-body #wraplone').html('');
        }
        $("#rmlOne").click(function() {
            $.post("/rmlabel", {
                    label: "1",
                    id: id
                },
                function(data, status) {
                    if (data) {
                        modal.find('.modal-body #wraplone').html('');
                        ep.find("span[name='detaillone']").html('None');
                    }
                });
        });
        if (ltwo != "None") {
            var remove = '<input type="button" value="Remove" class="btn btn-default btn-xs" id="rmlTwo"/>';
            modal.find('.modal-body #wrapltwo').html(ltwo + remove);
        } else {
            modal.find('.modal-body #wrapltwo').html('');
        }
        $("#rmlTwo").click(function() {
            $.post("/rmlabel", {
                    label: "2",
                    id: id
                },
                function(data, status) {
                    if (data) {
                        modal.find('.modal-body #wrapltwo').html('');
                        ep.find("span[name='detailltwo']").html('None');
                    }
                });
        });
    });
    $('#Modal-User').on('show.bs.modal', function(event) {
        var element = $(event.relatedTarget) // element that triggered the modal
        var ep = element.parent().parent();
        var id = ep.find("span[name='userid']").html();
        var nick = ep.find("span[name='usernick']").html();
        var host = ep.find("span[name='userhost']").html();
        var time = ep.find("span[name='userctime']").html();
        var by = ep.find("span[name='usercby']").html();
        var modal = $(this)
        $('#updatedid').val(id);
        modal.find('.modal-body #usernick').val(nick);
        modal.find('.modal-body #userhost').val(host);
        modal.find('.modal-body #createdtime').val(time);
        modal.find('.modal-body #createdby').val(by);
    });
    $(".catedel").click(function() {
        var ep = $(this).parent().parent();
        $.post("/category/del", {
                id: ep.find("span[name='cateid']").html()
            },
            function(data, status) {
                if (data) {
                    location.href = "/category?page=" + $("#pageIndex").val();
                } else {
                    alert("error");
                }
            });
    });
    $(".subcatedel").click(function() {
        var ep = $(this).parent().parent();
        $.post("/subcategory/del", {
                id: ep.find("span[name='subcateid']").html()
            },
            function(data, status) {
                if (data) {
                    location.href = "/subcategory?id="+$("select[name='category']").val()+"&page=" + $("#pageIndex").html();
                } else {
                    alert("error");
                }
            });
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
