$(document).ready(function() {
    $("#signupBtn").click(function(e){
        if($("#signupForm")[0].checkValidity()) {
            var formData = $("#signupForm").serializeArray();
            var formObj = {};

            $.map(formData, function(n, i){
                formObj[n['name']] = n['value'];
            });

            console.log(formObj);

            $.ajax({
                type: "POST",
                // PROD: url: "api.hiveapp.org/v1/subscribe_email",
                url: "v1/subscribe_email",
                data: JSON.stringify(formObj),
                success: function(){},
                dataType: "json",
                contentType : "application/json"
            });
            e.preventDefault();
        }
    }); 
});

