window.addEventListener('pageshow', function(event) {
    if(!event.persisted){
        collapseAll()
    }
});
document.getElementById('CategoryFilterForm').addEventListener('submit', async function(event) {

        event.preventDefault()
        const formData = new FormData(this);
  
        const response = await fetch('/main/tenders', {
            method: 'POST',
            body: formData
        });
        if(response.ok){
            const fullHtml = await response.text();
            
            const tempDiv = document.createElement('div');
            tempDiv.innerHTML = fullHtml;
            
            const newElement = tempDiv.querySelector('#tendersList');
            
            if (newElement) {
                const currentElement = document.querySelector('#tendersList');
                
                if (currentElement) {
                    currentElement.replaceWith(newElement);
                }
            }
        }
    });
function applySort(sortValue) {
    const tendersList = document.getElementById('tendersList');
    tendersList.style.opacity = '0.5';
    
    const url = new URL(window.location);
    url.searchParams.set('sort', sortValue);
    
   window.location.href = url;
}
// Переключение одной категории
function toggleCategory(btn) {
    const categoryItem = btn.closest('.category-item');
    const children = categoryItem.querySelector(':scope > .category-children');
    
    if (children) {
        if (children.classList.contains('collapsed')) {
            children.classList.remove('collapsed');
            btn.textContent = '▼';
        } else {
            children.classList.add('collapsed');
            btn.textContent = '▶';
        }
    }
}

// Развернуть все категории
function expandAll() {
    document.querySelectorAll('.category-children').forEach(el => {
        el.classList.remove('collapsed');
    });
    document.querySelectorAll('.toggle-btn').forEach(btn => {
        btn.textContent = '▼';
    });
}

// Свернуть все категории (кроме корневых, если нужно)
function collapseAll() {
    document.querySelectorAll('.category-children').forEach(el => {
        el.classList.add('collapsed');
    });
    document.querySelectorAll('.toggle-btn').forEach(btn => {
        btn.textContent = '▶';
    });
}

