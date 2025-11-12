// Create Report Page

document.addEventListener('DOMContentLoaded', () => {
    // Require authentication
    if (!requireAuth()) return;

    const reportForm = document.getElementById('reportForm');
    const beforeImageInput = document.getElementById('beforeImage');
    const imagePreview = document.getElementById('imagePreview');

    // Image preview
    if (beforeImageInput && imagePreview) {
        beforeImageInput.addEventListener('input', debounce(() => {
            const url = beforeImageInput.value.trim();
            if (url && isValidUrl(url)) {
                imagePreview.innerHTML = `<img src="${escapeHtml(url)}" alt="Preview" onerror="this.parentElement.classList.remove('show')">`;
                imagePreview.classList.add('show');
            } else {
                imagePreview.classList.remove('show');
            }
        }, 500));
    }

    // Form submission
    if (reportForm) {
        reportForm.addEventListener('submit', handleCreateReport);
    }

    // Mobile nav
    setupMobileNav();
});

async function handleCreateReport(e) {
    e.preventDefault();
    clearAllErrors();

    const submitBtn = document.getElementById('submitBtn');
    const title = document.getElementById('title').value.trim();
    const category = document.getElementById('category').value;
    const location = document.getElementById('location').value.trim();
    const description = document.getElementById('description').value.trim();
    const beforeImage = document.getElementById('beforeImage').value.trim();

    // Client-side validation
    let hasError = false;

    if (title.length === 0 || title.length > 200) {
        showFormError('title', 'Title must be between 1 and 200 characters');
        hasError = true;
    }

    if (!category) {
        showFormError('category', 'Please select a category');
        hasError = true;
    }

    if (location.length === 0 || location.length > 500) {
        showFormError('location', 'Location must be between 1 and 500 characters');
        hasError = true;
    }

    if (description.length === 0 || description.length > 2000) {
        showFormError('description', 'Description must be between 1 and 2000 characters');
        hasError = true;
    }

    if (!beforeImage || !isValidUrl(beforeImage)) {
        showFormError('beforeImage', 'Please provide a valid image URL');
        hasError = true;
    }

    if (hasError) return;

    setLoadingState(submitBtn, true);

    try {
        const data = await reportsApi.create({
            title,
            category,
            location,
            description,
            before_image: beforeImage
        });

        if (data.report) {
            showToast('Report created successfully!', 'success');
            setTimeout(() => {
                window.location.href = `report-detail.html?id=${data.report.id}`;
            }, 1000);
        }
    } catch (error) {
        console.error('Create report error:', error);
        handleApiError(error, {
            before_image: 'beforeImage'
        });
    } finally {
        setLoadingState(submitBtn, false);
    }
}

function setupMobileNav() {
    const navToggle = document.getElementById('navToggle');
    const navMenu = document.getElementById('navMenu');

    if (navToggle && navMenu) {
        navToggle.addEventListener('click', () => {
            navMenu.classList.toggle('active');
        });
    }
}
