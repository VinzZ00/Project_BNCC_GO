$(document).ready(() => {
    const editBtn = $("#edit-mode-btn");
    const saveBtn = $("#save-changes-btn");

    editBtn.click(() => {
        $(".memory-btns-container button").css("opacity", "100%");
        $("#edit-mode-btn").addClass("d-none");
        $("#save-changes-btn").removeClass("d-none");
        $("#back-btn").addClass("d-none");
    });

    saveBtn.click(() => {
        $(".memory-btns-container button").css("opacity", "0%");
        $("#edit-mode-btn").removeClass("d-none");
        $("#save-changes-btn").addClass("d-none");
        $("#back-btn").removeClass("d-none");
    });
});
