// notifications
function successAlert(x) {
    notie.alert({
        type: 'success',
        text: x,
    })
}

function errorAlert(x) {
    notie.alert({
        type: 'error',
        text: x,
    })
}

function warningAlert(x) {
    notie.alert({
        type: 'warning',
        text: x,
    })
}

// attention
function Prompt() {
    function toast(c) {
        const {
            msg = "",
            icon = "success",
            timer = 3000,
            showCloseButton = false,
        } = c;
        const Toast = Swal.mixin({
            toast: true,
            position: 'top',
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            onOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({
            icon: icon,
            title: msg,
            timer: timer,
            showCloseButton: showCloseButton,
        })
    }

    function alert(c) {
        const {
            cancelButton = false,
            html = "Some warning",
            icon = "success",
            confirmButtonText = "OK",
            confirmButton = true,
            title = '',
        } = c;
        const {value: result} = Swal.fire({
            html: html,
            confirmButtonText: confirmButtonText,
            icon: icon,
            showCancelButton: cancelButton,
            showConfirmButton: confirmButton,
            title: title,
        }).then((result) => {
            if (result) {
                if (result.dismiss === Swal.DismissReason.cancel) {
                } else {

                }
            }
        })
    }

    function confirm(c) {
        const {
            cancelButton = true,
            html = "Are you sure?",
            icon = "warning",
            confirmButtonText = "OK",
        } = c;
        const {value: result} = Swal.fire({
            html: html,
            confirmButtonText: confirmButtonText,
            icon: icon,
            showCancelButton: cancelButton,
            backdrop: false,
            showCloseButton: false,
        }).then((result) => {
            if (result) {
                if (result.dismiss !== Swal.DismissReason.cancel) {
                    if (c.callback !== undefined) {
                        c.callback(true)
                    }
                } else {
                    c.callback(false)
                }
            } else {
                c.callback(false)
            }
        })
    }

    function promptConfirm(c) {
        const {
            cancelButton = true,
            html = "Are you sure?",
            icon = "warning",
            confirmButtonText = "OK",
            confirmationText = "Yes",
        } = c;

        const {value: result} = Swal.fire({
            html: html,
            confirmButtonText: confirmButtonText,
            icon: icon,
            input: 'text',
            showCancelButton: cancelButton,
            backdrop: false,
            showCloseButton: false,
        }).then((result) => {
            if (result) {
                if (result.dismiss !== Swal.DismissReason.cancel) {
                    if (result.value.toUpperCase() === confirmationText.toUpperCase()) {
                        if (c.callback !== undefined) {
                            c.callback(true);
                        }
                    } else {
                        c.callback(false);
                    }
                } else {
                    c.callback(false);
                }
            }
        })
    }

    function prompt(c) {
        const {
            cancelButton = true,
            html = "Enter a value",
            icon = "question",
            buttonName = "OK",
            input = 'text',
            inputValue = '',
        } = c;

        const {value: result} = Swal.fire({
            html: html,
            confirmButtonText: buttonName,
            icon: icon,
            input: input,
            inputValue: inputValue,
            showCancelButton: cancelButton,
        }).then((result) => {
            if (result) {
                if (result.dismiss !== Swal.DismissReason.cancel) {
                    if (result.value !== "") {
                        if (c.callback !== undefined) {
                            c.callback(result);
                        }
                    } else {
                        c.callback(false);
                    }
                } else {
                    c.callback(false);
                }
            }
        })
    }

    return {
        confirm: confirm,
        alert: alert,
        promptConfirm: promptConfirm,
        prompt: prompt,
        toast: toast,
    };
}