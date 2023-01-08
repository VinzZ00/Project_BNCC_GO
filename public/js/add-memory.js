const resizeImage = (input, container, images) => {
    const uploadedImages = input.prop("files");

    for (let image of uploadedImages) {
        const newImage = new Image();
        const reader = new FileReader();

        newImage.addEventListener("load", () => {
            const width = newImage.naturalWidth;
            const height = newImage.naturalHeight;
            const aspectRatio = width / height;

            let newWidth = width,
                newHeight = height;

            if (aspectRatio > 1) {
                newWidth = height;
            } else if (aspectRatio < 1) {
                newHeight = width;
            }

            const canvas = document.createElement("canvas");
            const finalImage = new Image();

            canvas.width = newWidth;
            canvas.height = newHeight;
            canvas.getContext("2d").drawImage(newImage, 0, 0);

            finalImage.src = canvas.toDataURL();
            finalImage.className = "img-fluid rounded mb-3";

            container.append(
                $("<div class='w-100 col-md-4'></div>").append(finalImage)
            );
        });

        reader.addEventListener("load", () => {
            images.push(
                reader.result.replace(/^data:image\/[a-z]+;base64,/, "")
            );
        });

        newImage.src = URL.createObjectURL(image);
        reader.readAsDataURL(image);
    }
};

const convertImagesToBase64 = (input, images) => {
    const uploadedImages = input.prop("files");

    for (let image of uploadedImages) {
        const reader = new FileReader();

        reader.addEventListener("load", () => {
            images.push(
                reader.result.replace(/^data:image\/[a-z]+;base64,/, "")
            );
        });

        reader.readAsDataURL(image);
    }
};

$(document).ready(() => {
    const fileInput = $("input[type='file']");
    const imagesContainer = $("#images-container > .row");
    const images = [];

    fileInput.change(() => {
        resizeImage(fileInput, imagesContainer, images);
        convertImagesToBase64(fileInput, images);
    });

    $("#submit-btn").click((e) => {
        e.preventDefault();
        $.ajax({
            method: "POST",
            url: "/api/memories",
            data: JSON.stringify({
                description: $("#memory-name").val(),
                pictures: images,
                tags: $("#memory-tags").val().split(","),
            }),
            processData: false,
            contentType: "application/json",
        })
            .done(() => {
                alert("Memory successfully created!");
                window.location.href = "/memories";
            })
            .fail((data) => {
                console.log(data);
                alert(data);
            });
    });
});
