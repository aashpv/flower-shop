document.addEventListener('DOMContentLoaded', function () {
    const deleteButton = document.getElementById('delete');

    const pathname = window.location.pathname;
    const pathParts = pathname.split('/');
    const id = pathParts[pathParts.length - 2]; // Get the ID from the URL

    const token = localStorage.getItem("Bearer");

    if (id) {
        fetch('/product/' + id)
            .then(response => {
                return response.json();
            })
            .then(data => {
                // Use the JSON data to populate the page
                const productDiv = document.createElement('div');
                productDiv.innerHTML = `
                    <h1>${data.product.name}</h1>
                    <p>Описание:</p>
                    <p>${data.product.description}</p>
                    <p>Цена: ${data.product.price}</p>
                `;
                document.getElementById('product').appendChild(productDiv);

                // Add event listener for delete button
                deleteButton.addEventListener('click', () => {
                    if (confirm('Вы действительно желаете удалить этот товар?')) {
                        fetch('/product/' + id, {
                            method: 'DELETE',
                            headers: { 'Authorization': 'Bearer ' + token }
                        })
                            .then(response => {
                                if (response.ok) {
                                    alert('Товар удален');
                                    window.location.href = '/';
                                } else {
                                    alert('Ошибка при удалении товара');
                                }
                            })
                            .catch(error => console.error('Error during delete:', error));
                    }
                });
            })
            .catch(error => console.error('Fetch error:', error));
    } else {
        console.error('Product ID not found in URL');
    }
});
