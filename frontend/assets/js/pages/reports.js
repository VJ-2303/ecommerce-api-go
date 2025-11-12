// Reports Page

let currentPage = 0;
const pageSize = 12;
let currentFilters = {
    status: '',
    category: ''
};

document.addEventListener('DOMContentLoaded', () => {
    setupMobileNav();
    setupFilters();
    setupPagination();
    loadReports();

    // Check for query parameters
    const statusParam = getQueryParam('status');
    const categoryParam = getQueryParam('category');

    if (statusParam) {
        document.getElementById('statusFilter').value = statusParam;
        currentFilters.status = statusParam;
    }

    if (categoryParam) {
        document.getElementById('categoryFilter').value = categoryParam;
        currentFilters.category = categoryParam;
    }
});

function setupFilters() {
    const statusFilter = document.getElementById('statusFilter');
    const categoryFilter = document.getElementById('categoryFilter');
    const clearFilters = document.getElementById('clearFilters');

    if (statusFilter) {
        statusFilter.addEventListener('change', (e) => {
            currentFilters.status = e.target.value;
            currentPage = 0;
            setQueryParam('status', e.target.value);
            loadReports();
        });
    }

    if (categoryFilter) {
        categoryFilter.addEventListener('change', (e) => {
            currentFilters.category = e.target.value;
            currentPage = 0;
            setQueryParam('category', e.target.value);
            loadReports();
        });
    }

    if (clearFilters) {
        clearFilters.addEventListener('click', () => {
            statusFilter.value = '';
            categoryFilter.value = '';
            currentFilters = { status: '', category: '' };
            currentPage = 0;
            window.history.pushState({}, '', window.location.pathname);
            loadReports();
        });
    }
}

function setupPagination() {
    const prevBtn = document.getElementById('prevBtn');
    const nextBtn = document.getElementById('nextBtn');

    if (prevBtn) {
        prevBtn.addEventListener('click', () => {
            if (currentPage > 0) {
                currentPage--;
                loadReports();
            }
        });
    }

    if (nextBtn) {
        nextBtn.addEventListener('click', () => {
            currentPage++;
            loadReports();
        });
    }
}

async function loadReports() {
    const grid = document.getElementById('reportsGrid');
    if (!grid) return;

    grid.innerHTML = '<div class="loading">Loading reports...</div>';

    try {
        const params = {
            limit: pageSize,
            offset: currentPage * pageSize,
            status: currentFilters.status,
            category: currentFilters.category
        };

        const data = await reportsApi.getAll(params);

        if (data.reports && data.reports.length > 0) {
            grid.innerHTML = '';
            data.reports.forEach(report => {
                grid.appendChild(createReportCard(report));
            });

            updatePagination(data.reports.length);
        } else {
            grid.innerHTML = '<div class="loading">No reports found</div>';
            updatePagination(0);
        }
    } catch (error) {
        console.error('Error loading reports:', error);
        grid.innerHTML = '<div class="loading">Failed to load reports</div>';
    }
}

function createReportCard(report) {
    const card = document.createElement('div');
    card.className = 'report-card';
    card.onclick = () => {
        window.location.href = `report-detail.html?id=${report.id}`;
    };

    const categoryIcon = CATEGORY_ICONS[report.category] || '📝';
    const statusClass = STATUS_CLASSES[report.status] || 'pending';

    card.innerHTML = `
        <img src="${escapeHtml(report.before_image)}" 
             alt="${escapeHtml(report.title)}" 
             class="report-image"
             onerror="this.src='https://via.placeholder.com/400x200?text=No+Image'">
        <div class="report-content">
            <div class="report-header">
                <h3 class="report-title">${escapeHtml(report.title)}</h3>
                <span class="report-status ${statusClass}">
                    ${report.status.replace('-', ' ')}
                </span>
            </div>
            <div class="report-category">
                <span>${categoryIcon}</span>
                <span>${report.category}</span>
            </div>
            <p class="report-description">${escapeHtml(report.description)}</p>
            <div class="report-footer">
                <div class="report-location">
                    <span>📍</span>
                    <span>${truncateText(escapeHtml(report.location), 30)}</span>
                </div>
                <div class="report-date">
                    <span>🕒</span>
                    <span>${formatDateOnly(report.created_at)}</span>
                </div>
            </div>
        </div>
    `;

    return card;
}

function updatePagination(itemsCount) {
    const prevBtn = document.getElementById('prevBtn');
    const nextBtn = document.getElementById('nextBtn');
    const paginationInfo = document.getElementById('paginationInfo');

    if (prevBtn) {
        prevBtn.disabled = currentPage === 0;
    }

    if (nextBtn) {
        nextBtn.disabled = itemsCount < pageSize;
    }

    if (paginationInfo) {
        paginationInfo.textContent = `Page ${currentPage + 1}`;
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
