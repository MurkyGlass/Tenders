window.addEventListener('pageshow', function(event) {
    if(event.persisted){
        window.location.href = this.window.location.pathname
    }
});
async function EditOfferForm(id) {
    const formData = new FormData(document.getElementById('offerForm'));
  
        const response = await fetch('/protected/lk/offers/'+id+'/edit', {//
            method: 'POST',
            body: formData
        });
        if (response.ok){
            alert("Предложение изменено")
            window.location.href = '/protected/lk/offers/'+id
        }else{
            alert("Ошибка:"+ await response.text())
        }
}
async function EditDraftOfferForm(id) {
    const formData = new FormData(document.getElementById('offerForm'));
  
        const response = await fetch('/protected/lk/offers/'+id+'/edit/draft', {//
            method: 'POST',
            body: formData
        });
        if (response.ok){
            alert("Предложение изменено")
            window.location.href = '/protected/lk/offers/'+id
        }else{
            alert("Ошибка:"+ await response.text())
        }
}
// todo validation for file format and size
document.getElementById('file_select').addEventListener('click', function(e) {
    // Сохраняем текущие файлы ДО открытия диалога
    const oldFiles = Array.from(this.files);
    console.log('Сохранено до выбора:', oldFiles.length);
    
    const input = this;
    
    const handleChange = () => {
        const newFiles = Array.from(input.files);
        console.log('Новые файлы:', newFiles.length);
        
        const dt = new DataTransfer();
        
        oldFiles.forEach(file => dt.items.add(file));
        
        newFiles.forEach(newFile => {
            const exists = oldFiles.some(oldFile => 
                oldFile.name === newFile.name && oldFile.size === newFile.size
            );
            
            if (!exists) {
                dt.items.add(newFile);
            }
        });
        
        input.files = dt.files;
        
        updateFileDisplay();
        
        input.removeEventListener('change', handleChange);
    };
    
    this.addEventListener('change', handleChange);
});

function removeFile(index) {
    const input = document.getElementById('file_select');
    
    const dt = new DataTransfer();
    
    for (let i = 0; i < input.files.length; i++) {
        if (i !== index) {
            dt.items.add(input.files[i]);
        }
    }
    
    input.files = dt.files;
    
    updateFileDisplay();
}

function updateFileDisplay() {
    const input = document.getElementById('file_select');
    const display = document.getElementById('file-name-display');
    const list = document.getElementById('file-list');
    
    if (input.files.length === 0) {
        display.textContent = 'Файлы не выбраны';
        list.innerHTML = '';
        return;
    }
    
    display.textContent = `Выбрано файлов: ${input.files.length}`;

    list.innerHTML = Array.from(input.files).map((file, index) => `
        <div class="file-list-item">
            <span>
                <i class="fas fa-file"></i> 
                ${file.name}
                <span class="file-size">(${(file.size / 1024).toFixed(1)} KB)</span>
            </span>
            <span class="remove-file" onclick="removeFile(${index})">
                <i class="fas fa-times"></i>
            </span>
        </div>
    `).join('');
}