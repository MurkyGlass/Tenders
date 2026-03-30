window.addEventListener('pageshow', function(event) {
    if(!event.persisted){
        collapseAll()
    }
});
async function Search(){
    searcher = document.getElementById('searchbox')
    sorter = document.getElementById('sortSelect')
    const url = new URL(window.location);
    url.searchParams.set('sort', sorter.value);
    url.searchParams.set('search',searcher.value)
    form = document.getElementById('CategoryFilterForm')
    
    const formData = new FormData(form);
    await UpdateTenders(url,formData)
        
}
async function UpdateTenders(url,formData) {
    const response = await fetch(url, {
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
}
async function createTenderLink(){
    const resp = await fetch("/protected/tender/create");
            if(!resp.ok){
                if(resp.status == 401){
                    document.getElementById('loginModal').classList.add("show")
                }else if(resp.status == 403){
                    await Refresh()
                    createTenderLink()
                }
                return
            } 
                 
            window.location.href = "/protected/tender/create";
}

document.getElementById('CategoryFilterForm').addEventListener('submit', async function(event) {

        searcher = document.getElementById('searchbox')
        sorter = document.getElementById('sortSelect')
        const url = new URL(window.location);
        url.searchParams.set('sort', sorter.value);
        url.searchParams.set('search',searcher.value)

        event.preventDefault()
        const formData = new FormData(this);
  
        await UpdateTenders(url,formData)
    });
async function applySort(sortValue) {
    searcher = document.getElementById('searchbox')
    const url = new URL(window.location);
    url.searchParams.set('sort', sortValue);
    url.searchParams.set('search',searcher.value)
    
    form = document.getElementById('CategoryFilterForm')
    
    const formData = new FormData(form);
  
    await UpdateTenders(url,formData)
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

