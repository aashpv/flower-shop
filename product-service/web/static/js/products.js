// document.addEventListener('DOMContentLoaded', () => {
//     const createProductBtn = document.getElementById('create-product');
//
//     // Получение токена из localStorage
//     const token = localStorage.getItem("Bearer")
//     // alert("token from local storage: " + token)
//     // Функция для отправки запроса на защищенный роут с токеном в заголовке
//     const sendAuthenticatedRequest = (url) => {
//         fetch(url, {
//             method: 'GET',
//             headers: {
//                 'Authorization': 'Bearer ' + token
//             }
//         })
//             .then(response => {
//                 if (response.ok) {
//                     // Если запрос успешен, перенаправляем пользователя на указанный URL
//                     window.location.href = url;
//                 } else {
//                     // В случае ошибки обрабатываем ее
//                     console.error('Ошибка при отправке запроса:', response.statusText);
//                 }
//             })
//             .catch(error => {
//                 console.error('Ошибка при отправке запроса:', error);
//             });
//     };
//
//     createProductBtn.addEventListener('click', () => {
//         // Отправка запроса на страницу создания нового товара с токеном в заголовке
//         sendAuthenticatedRequest('/create');
//     });
//
//     fetch('/products')
//         .then(response => {
//             return response.json()
//         })
//         .then(data => {
//             const productsDiv = document.getElementById('products');
//             const products = data.product;
//             products.forEach(product => {
//                 const productDiv = document.createElement('div');
//                 productDiv.innerHTML = `
//         <h2>${product.name}</h2>
//         <p>Описание:</p>
//         <p>${product.description}</p>
//         <p>Цена: ${product.price}</p>
//       `;
//                 productsDiv.appendChild(productDiv);
//             });
//         })
//         .catch(error => console.error(error));
// });

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
                        }
                    })
                    .catch(error => {
                        console.error('Ошибка при отправке запроса:', error);
                    })
            })
        })
        .catch(error => console.error(error));
});