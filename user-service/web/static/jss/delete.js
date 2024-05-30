document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#del')

    const token = localStorage.getItem('Bearer')

    form.addEventListener('submit', async (event) => {
        event.preventDefault()

        const email = document.querySelector('#email').value

        const data = {
            email,
        }

        const jsonData = JSON.stringify(data)

        const response = await fetch("/delete", {
            method: "DELETE",
            headers: {
                "Authorization": 'Bearer ' + token,
                "Content-Type": 'application/json'
            },
            body: jsonData
        })

        if (response.ok){
            const data = await response.json()

            if (data.status === 200) {
                alert("Данные доставлены")
                window.location.href = '/'
            }else{
                alert('Ошибка при отправке данных');
            }
        }else{
            alert('Ошибка при отправке данных');
        }
    })
})