$(document).ready(function() {
    $("#signupBtn").click(function(e){
        if($("#signupForm")[0].checkValidity()) {
            var formData = $("#signupForm").serializeArray();
            var formObj = {};

            $.map(formData, function(n, i){
                if (n['name'] === 'classYear') {
                    formObj[n['name']] = parseInt(n['value']);
                } else {
                    formObj[n['name']] = n['value'];
                }
            });

            $.ajax({
                method: "POST",
                url: "https://api.hiveapp.org/v1/subscribe_email",
                headers: {
                    'Access-Control-Allow-Origin': '*'
                },
                data: JSON.stringify(formObj),
                success: function(){},
                dataType: "json",
                contentType : "application/json"
            });
            e.preventDefault();
        }
    });
});
