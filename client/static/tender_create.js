/**
 * Скрипт для страницы создания тендера
 * Обрабатывает валидацию форм и взаимодействия с пользователем
 */

document.addEventListener('DOMContentLoaded', function() {
    'use strict';
    
    // Получаем элементы формы
    const tenderForm = document.getElementById('tenderForm');
    const startInput = document.getElementById('datetime_start');
    const endInput = document.getElementById('datetime_end');
    const warning = document.getElementById('dateWarning');
    const saveDraftBtn = document.getElementById('saveDraftBtn');
    const statusSelect = document.getElementById('id_status');
    
    // Инициализация
    initDateValidation();
    initFormValidation();
    initDraftButton();
    initAutoResizeTextarea();
    
    /**
     * Инициализация валидации дат
     */
    function initDateValidation() {
        if (!startInput || !endInput || !warning) return;
        
        function validateDates() {
            if (startInput.value && endInput.value) {
                const start = new Date(startInput.value);
                const end = new Date(endInput.value);
                const now = new Date();
                
                // Проверка, что дата начала не в прошлом
                if (start < now) {
                    showFieldError(startInput, 'Дата начала не может быть в прошлом');
                } else {
                    clearFieldError(startInput);
                }
                
                // Проверка, что дата окончания позже даты начала
                if (end <= start) {
                    warning.classList.add('show');
                    endInput.setCustomValidity('Дата окончания должна быть позже даты начала');
                    showFieldError(endInput, 'Дата окончания должна быть позже даты начала');
                } else {
                    warning.classList.remove('show');
                    endInput.setCustomValidity('');
                    clearFieldError(endInput);
                }
            }
        }
        
        startInput.addEventListener('change', validateDates);
        endInput.addEventListener('change', validateDates);
        
        // Валидация при вводе
        startInput.addEventListener('input', validateDates);
        endInput.addEventListener('input', validateDates);
    }
    
    /**
     * Инициализация валидации формы
     */
    function initFormValidation() {
        if (!tenderForm) return;
        
        tenderForm.addEventListener('submit', function(e) {
            let isValid = true;
            
            // Валидация названия
            const nameInput = document.getElementById('name');
            if (nameInput && nameInput.value.trim().length < 3) {
                showFieldError(nameInput, 'Название должно содержать минимум 3 символа');
                isValid = false;
            }
            
            // Валидация дат
            if (startInput && endInput) {
                const start = new Date(startInput.value);
                const end = new Date(endInput.value);
                const now = new Date();
                
                if (start < now) {
                    showFieldError(startInput, 'Дата начала не может быть в прошлом');
                    isValid = false;
                }
                
                if (end <= start) {
                    showFieldError(endInput, 'Дата окончания должна быть позже даты начала');
                    isValid = false;
                }
            }
            
            // Валидация описания (если есть)
            const descInput = document.getElementById('description');
            if (descInput && descInput.value.length > 500) {
                showFieldError(descInput, 'Описание не может превышать 500 символов');
                isValid = false;
            }
            
            if (!isValid) {
                e.preventDefault();
                scrollToFirstError();
            }
        });
    }
    
    /**
     * Инициализация кнопки сохранения черновика
     */
    function initDraftButton() {
        if (!saveDraftBtn || !statusSelect) return;
        
        saveDraftBtn.addEventListener('click', function() {
            // Устанавливаем статус "Черновик" (id = 1)
            if (!statusSelect.value) {
                statusSelect.value = '1';
            }
            
            // Показываем уведомление
            showNotification('Тендер сохранен как черновик', 'info');
            
            // Отправляем форму
            if (tenderForm) {
                tenderForm.submit();
            }
        });
    }
    
    /**
     * Автоматическое расширение textarea
     */
    function initAutoResizeTextarea() {
        const textarea = document.getElementById('description');
        if (!textarea) return;
        
        textarea.addEventListener('input', function() {
            this.style.height = 'auto';
            this.style.height = (this.scrollHeight) + 'px';
            
            // Ограничиваем максимальную высоту
            if (this.scrollHeight > 200) {
                this.style.height = '200px';
                this.style.overflowY = 'auto';
            }
        });
    }
    
    /**
     * Показать ошибку поля
     */
    function showFieldError(field, message) {
        if (!field) return;
        
        field.classList.add('error');
        
        // Удаляем существующее сообщение об ошибке
        const existingError = field.parentNode.querySelector('.error-message');
        if (existingError) {
            existingError.remove();
        }
        
        // Создаем новое сообщение
        const errorDiv = document.createElement('div');
        errorDiv.className = 'error-message';
        errorDiv.style.color = '#dc3545';
        errorDiv.style.fontSize = '0.85rem';
        errorDiv.style.marginTop = '5px';
        errorDiv.innerHTML = `<i class="fas fa-exclamation-circle"></i> ${message}`;
        
        field.parentNode.appendChild(errorDiv);
    }
    
    /**
     * Очистить ошибку поля
     */
    function clearFieldError(field) {
        if (!field) return;
        
        field.classList.remove('error');
        const errorMsg = field.parentNode.querySelector('.error-message');
        if (errorMsg) {
            errorMsg.remove();
        }
    }
    
    /**
     * Прокрутка к первому полю с ошибкой
     */
    function scrollToFirstError() {
        const firstError = document.querySelector('.form-input.error');
        if (firstError) {
            firstError.scrollIntoView({ behavior: 'smooth', block: 'center' });
            firstError.focus();
        }
    }
    
    /**
     * Показать уведомление
     */
    function showNotification(message, type = 'info') {
        // Создаем элемент уведомления
        const notification = document.createElement('div');
        notification.className = `notification notification-${type}`;
        notification.innerHTML = `
            <i class="fas ${type === 'success' ? 'fa-check-circle' : 'fa-info-circle'}"></i>
            <span>${message}</span>
        `;
        
        // Стили для уведомления
        notification.style.position = 'fixed';
        notification.style.top = '20px';
        notification.style.right = '20px';
        notification.style.padding = '15px 25px';
        notification.style.background = type === 'success' ? '#28a745' : '#17a2b8';
        notification.style.color = 'white';
        notification.style.borderRadius = '8px';
        notification.style.boxShadow = '0 4px 12px rgba(0,0,0,0.15)';
        notification.style.zIndex = '9999';
        notification.style.display = 'flex';
        notification.style.alignItems = 'center';
        notification.style.gap = '10px';
        notification.style.animation = 'slideIn 0.3s ease';
        
        document.body.appendChild(notification);
        
        // Удаляем через 3 секунды
        setTimeout(() => {
            notification.style.animation = 'slideOut 0.3s ease';
            setTimeout(() => {
                notification.remove();
            }, 300);
        }, 3000);
    }
    
    // Добавляем анимации для уведомлений
    const style = document.createElement('style');
    style.textContent = `
        @keyframes slideIn {
            from {
                transform: translateX(100%);
                opacity: 0;
            }
            to {
                transform: translateX(0);
                opacity: 1;
            }
        }
        
        @keyframes slideOut {
            from {
                transform: translateX(0);
                opacity: 1;
            }
            to {
                transform: translateX(100%);
                opacity: 0;
            }
        }
    `;
    document.head.appendChild(style);
    
    // Предотвращаем отправку формы с невалидными датами
    if (tenderForm) {
        tenderForm.addEventListener('submit', function(e) {
            if (startInput && endInput) {
                const start = new Date(startInput.value);
                const end = new Date(endInput.value);
                
                if (end <= start) {
                    e.preventDefault();
                    showNotification('Пожалуйста, исправьте ошибки в датах', 'error');
                }
            }
        });
    }
});