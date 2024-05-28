document.addEventListener('DOMContentLoaded', () => {
    const createProductBtn = document.querySelector('#create-product');

    // Получение токена из localStorage
    const token = localStorage.getItem("Bearer")
    alert("token from local storage: " + token)
    // Функция для отправки запроса на защищенный роут с токеном в заголовке
    const sendAuthenticatedRequest = (url) => {
        fetch(url, {
            method: 'GET',
            headers: {
                'Authorization': 'Bearer ' + token
            }
        })
            .then(response => {
                if (response.ok) {
                    // Если запрос успешен, перенаправляем пользователя на указанный URL
                    window.location.href = url;
                } else {
                    // В случае ошибки обрабатываем ее
                    console.error('Ошибка при отправке запроса:', response.statusText);
                }
            })
            .catch(error => {
                console.error('Ошибка при отправке запроса:', error);
            });
    };

    createProductBtn.addEventListener('click', () => {
        // Отправка запроса на страницу создания нового товара с токеном в заголовке
        sendAuthenticatedRequest('/create');
    });

    fetch('/products')
        .then(response => response.json())
        .then(data => {
            const productsDiv = document.getElementById('products');
            const products = data.products;
            products.forEach(product => {
                const productDiv = document.createElement('div');
                productDiv.innerHTML = `
        <h2>${product.Name}</h2>
        <p>Описание: ${product.Description}</p>
        <p>Цена: ${product.Price}</p>
      `;
                productsDiv.appendChild(productDiv);
            });
        })
        .catch(error => console.error(error));
});