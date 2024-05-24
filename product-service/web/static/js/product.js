const pathname = window.location.pathname;
const pathParts = pathname.split('/');
const id = pathParts[pathParts.length - 2]; // Get the ID from the URL
const token = localStorage.getItem("Bearer")

if (id) {
    fetch('/product/' + id)
        .then(response => response.json())
        .then(data => {
            // Use the JSON data to populate the page
            const productDiv = document.createElement('div');
            productDiv.innerHTML = `
        <h1>${data.product.Name}</h1>
        <p>Описание: ${data.product.Description}</p>
        <p>Цена: ${data.product.Price}</p>
      `;
            document.getElementById('product').appendChild(productDiv);

            // Add event listener for delete button
            const deleteButton = document.getElementById('delete');
            deleteButton.addEventListener('click', () => {
                if (confirm('Вы действительно желаете удалить этот товар?')) {
                    fetch('/product/' + id, { method: 'DELETE' , headers: {'Authorization': 'Bearer '+ token}})
                        .then(response => {
                            if (response.ok) {
                                alert('Товар удален');
                                window.location.href = '/';
                            } else {
                                alert('Ошибка при удалении товара');
                            }
                        })
                        .catch(error => console.error(error));
                }
            });
        })
        .catch(error => console.error(error));
} else {
    console.error('Product ID not found in URL');
}