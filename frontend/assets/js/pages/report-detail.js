// Report Detail Page

let reportId = null;
let currentReport = null;

document.addEventListener('DOMContentLoaded', () => {
    reportId = getQueryParam('id');

    if (!reportId) {
        showToast('Report not found', 'error');
        setTimeout(() => {
            window.location.href = 'reports.html';
        }, 1500);
        return;
    }

    loadReportDetail();

    // Setup update form if admin
    if (isAdmin()) {
        setupAdminPanel();
    }
});

async function loadReportDetail() {
    const container = document.getElementById('reportDetail');
    if (!container) return;

    container.innerHTML = '<div class="loading">Loading report details...</div>';

    try {
        const data = await reportsApi.getById(reportId);

        if (data.report) {
            currentReport = data.report;
            displayReportDetail(data.report);
            
            if (isAdmin()) {
                showAdminPanel(data.report);
            }
        }
    } catch (error) {
        console.error('Error loading report:', error);
        container.innerHTML = '<div class="loading">Failed to load report</div>';
        
        if (error.error === 'this resourse is not found') {
            showToast('Report not found', 'error');
            setTimeout(() => {
                window.location.href = 'reports.html';
            }, 1500);
        }
    }
}

function displayReportDetail(report) {
    const container = document.getElementById('reportDetail');
    const categoryIcon = CATEGORY_ICONS[report.category] || '📝';
    const statusClass = STATUS_CLASSES[report.status] || 'pending';

    const hasAfterImage = report.after_image && report.after_image !== '';

    container.innerHTML = `
        <div class="detail-image-container">
            <div class="detail-image-wrapper">
                <div class="detail-image-label">Before</div>
                <img src="${escapeHtml(report.before_image)}" 
                     alt="Before" 
                     class="detail-image"
                     onerror="this.src='https://via.placeholder.com/600x300?text=No+Image'">
            </div>
            ${hasAfterImage ? `
                <div class="detail-image-wrapper">
                    <div class="detail-image-label">After</div>
                    <img src="${escapeHtml(report.after_image)}" 
                         alt="After" 
                         class="detail-image"
                         onerror="this.src='https://via.placeholder.com/600x300?text=No+Image'">
                </div>
            ` : ''}
        </div>
        <div class="detail-content">
            <div style="display: flex; justify-content: space-between; align-items: start; margin-bottom: 1rem; flex-wrap: wrap; gap: 1rem;">
                <h1 style="margin: 0;">${escapeHtml(report.title)}</h1>
                <span class="report-status ${statusClass}" style="font-size: 1rem;">
                    ${report.status.replace('-', ' ')}
                </span>
            </div>
            
            <div class="detail-meta">
                <div class="report-category" style="font-size: 1rem;">
                    <span>${categoryIcon}</span>
                    <span>${report.category}</span>
                </div>
                <div style="color: var(--text-secondary);">
                    <span>📍</span>
                    <span>${escapeHtml(report.location)}</span>
                </div>
                <div style="color: var(--text-secondary);">
                    <span>👤</span>
                    <span>${escapeHtml(report.user_name || 'Anonymous')}</span>
                </div>
                <div style="color: var(--text-secondary);">
                    <span>🕒</span>
                    <span>${formatDate(report.created_at)}</span>
                </div>
            </div>

            <div style="margin-bottom: 1.5rem;">
                <h3 style="margin-bottom: 0.5rem;">Description</h3>
                <p style="color: var(--text-secondary); line-height: 1.6;">
                    ${escapeHtml(report.description)}
                </p>
            </div>

            ${report.completed_at ? `
                <div style="color: var(--text-secondary);">
                    <strong>Completed At:</strong> ${formatDate(report.completed_at)}
                </div>
            ` : ''}

            <div style="color: var(--text-secondary); font-size: 0.875rem; margin-top: 1.5rem; padding-top: 1rem; border-top: 1px solid var(--border-color);">
                <div><strong>Report ID:</strong> ${report.id}</div>
                <div><strong>Last Updated:</strong> ${formatDate(report.updated_at)}</div>
            </div>
        </div>
    `;
}

function setupAdminPanel() {
    const updateForm = document.getElementById('updateForm');
    if (updateForm) {
        updateForm.addEventListener('submit', handleUpdateReport);
    }
}

function showAdminPanel(report) {
    const adminPanel = document.getElementById('adminPanel');
    if (!adminPanel) return;

    adminPanel.classList.remove('hidden');

    // Pre-fill current status
    const statusSelect = document.getElementById('status');
    if (statusSelect) {
        statusSelect.value = report.status;
    }

    // Pre-fill after image if exists
    const afterImageInput = document.getElementById('afterImage');
    if (afterImageInput && report.after_image) {
        afterImageInput.value = report.after_image;
    }
}

async function handleUpdateReport(e) {
    e.preventDefault();
    
    const updateBtn = document.getElementById('updateBtn');
    const updateError = document.getElementById('updateError');
    const status = document.getElementById('status').value;
    const afterImage = document.getElementById('afterImage').value.trim();

    // Clear previous errors
    updateError.textContent = '';
    updateError.classList.remove('show');

    // Validate
    if (status === 'completed' && !afterImage) {
        updateError.textContent = 'After image is required when marking as completed';
        updateError.classList.add('show');
        return;
    }

    setLoadingState(updateBtn, true);

    try {
        await reportsApi.update(reportId, {
            status,
            after_image: afterImage
        });

        showToast('Report updated successfully!', 'success');
        
        // Reload report detail
        setTimeout(() => {
            loadReportDetail();
        }, 1000);
    } catch (error) {
        console.error('Update error:', error);
        updateError.textContent = error.error || 'Failed to update report';
        updateError.classList.add('show');
    } finally {
        setLoadingState(updateBtn, false);
    }
}
