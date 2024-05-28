document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#signup');

    form.addEventListener('submit', async (event) => {
        event.preventDefault();

        const first_name = document.querySelector('#first_name').value;
        const last_name = document.querySelector('#last_name').value;
        const address = document.querySelector('#address').value;
        const phone = document.querySelector('#phone').value;
        const email = document.querySelector('#email').value;
        const password = document.querySelector('#password').value;

        const data = {
            first_name,
            last_name,
            address,
            phone,
            email,
            password,
        };

        const jsonData = JSON.stringify(data);

        const response = await fetch('/signup', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: jsonData,
        });

        if (response.ok) {
            const data = await response.json();
            if (data.status === 200) {
                alert('Данные доставлены');
                window.location.href = '/login'; // перенаправление на страницу входа
            } else {
                alert('Ошибка при отправке данных');
            }
        } else {
            alert('Ошибка при отправке данных');
        }
    });
});