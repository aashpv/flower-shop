document.addEventListener('DOMContentLoaded', () => {
    const createProductBtn = document.getElementById('create-product');

    const token = localStorage.getItem("Bearer")

    fetch('/products')
        .then(response => {
            return response.json()
        })
        .then(data => {
            data.products.forEach(product => {
                const productDiv = document.createElement('div');
                productDiv.innerHTML = `
                <h1>${product.name}</h1>
                <p>Описание:</p>
                <p>${product.description}</p>
                <p>Цена: ${product.price}</p>
              `;
                document.getElementById('products').appendChild(productDiv);
            });

            // Add event listener for create button
            createProductBtn.addEventListener('click', () => {
                fetch('/create', {
                    method: 'GET',
                    headers: { 'Authorization': 'Bearer ' + token }
                })
                    .then(response => {
                        if (response.ok) {
                            // Если запрос успешен, перенаправляем пользователя на указанный URL
                            window.location.href = '/create';
                        } else {
                            // В случае ошибки обрабатываем ее
                            console.error('Ошибка при отправке запроса:', response.statusText);
                            alert("Не имеешь права!")
                        }
                    })
                    .catch(error => {
                        console.error('Ошибка при отправке запроса:', error);
                    })
            })
        })
        .catch(error => console.error(error));
});