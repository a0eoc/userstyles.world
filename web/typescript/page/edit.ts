export function editStyleFields() {
    const description = document.querySelector('#description');
    description && description.addEventListener('input', () => {
        var len = description.value.length;
        if(len >= 100) {
            console.log(len);
        }
    });
}
