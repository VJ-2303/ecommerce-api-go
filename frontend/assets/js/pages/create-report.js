// Create Report Page

document.addEventListener('DOMContentLoaded', () => {
    // Require authentication
    if (!requireAuth()) return;

    const reportForm = document.getElementById('reportForm');
    const beforeImageInput = document.getElementById('beforeImage');
    const imagePreview = document.getElementById('imagePreview');

    // Image preview
    if (beforeImageInput && imagePreview) {
        beforeImageInput.addEventListener('change', (e) => {
            const file = e.target.files[0];
            if (file) {
                // Validate file size (5MB max)
                if (file.size > 5 * 1024 * 1024) {
                    showFormError('beforeImage', 'Image size must be less than 5MB');
                    beforeImageInput.value = '';
                    imagePreview.classList.remove('show');
                    return;
                }

                // Validate file type
                if (!file.type.startsWith('image/')) {
                    showFormError('beforeImage', 'Please upload an image file');
                    beforeImageInput.value = '';
                    imagePreview.classList.remove('show');
                    return;
                }

                hideFormError('beforeImage');

                // Preview image
                const reader = new FileReader();
                reader.onload = (e) => {
                    imagePreview.innerHTML = `<img src="${e.target.result}" alt="Preview">`;
                    imagePreview.classList.add('show');
                };
                reader.readAsDataURL(file);
            } else {
                imagePreview.classList.remove('show');
            }
        });
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
    const beforeImageFile = document.getElementById('beforeImage').files[0];

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

    if (!beforeImageFile) {
        showFormError('beforeImage', 'Please upload an image');
        hasError = true;
    } else if (beforeImageFile.size > 5 * 1024 * 1024) {
        showFormError('beforeImage', 'Image size must be less than 5MB');
        hasError = true;
    }

    if (hasError) return;

    setLoadingState(submitBtn, true);

    try {
        // Convert image to base64
        const beforeImageBase64 = await fileToBase64(beforeImageFile);

        const data = await reportsApi.create({
            title,
            category,
            location,
            description,
            before_image: beforeImageBase64
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
