$(document).ready(function() {
    $("#signupBtn").click(function(){
        if($("#signupForm")[0].checkValidity()) {
            var formData = $("#signupForm").serializeArray();
            var formObj = {};

            $.map(formData, function(n, i){
                formObj[n['name']] = n['value'];
            });

            console.log(formObj);

            $.ajax({
                type: "POST",
                url: "v1/subscribe_email",
                data: JSON.stringify(formObj),
                success: function(){},
                dataType: "json",
                contentType : "application/json"
            });
        }
    }); 
});

