document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#create');

    const token = localStorage.getItem("Bearer")
    alert(token)

    form.addEventListener('submit', (event) => {
        event.preventDefault();

        const name = document.querySelector('#name').value;
        const description = document.querySelector('#description').value;
        const price = parseInt(document.querySelector('#price').value);

        const data = {
            name,
            description,
            price,
        };

        const jsonData = JSON.stringify(data);

        console.log(jsonData)

        fetch('/create', {
            method: 'POST',
            headers: {
                'Authorization': 'Bearer ' + token,
                'Content-Type': 'application/json'
            },
            body: jsonData,
        }).then((response) => {
            if (response.ok) {
                alert('Данные доставлены');
                window.location.href = '/'
            } else {
                alert('Ошибка при отправке данных');
            }
        });
    });
});