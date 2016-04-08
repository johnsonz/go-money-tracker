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
});
