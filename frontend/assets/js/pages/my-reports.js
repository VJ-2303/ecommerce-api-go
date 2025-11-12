// My Reports Page

let currentPage = 0;
const pageSize = 12;

document.addEventListener('DOMContentLoaded', () => {
    // Require authentication
    if (!requireAuth()) return;

    setupMobileNav();
    setupPagination();
    loadMyReports();
});

function setupPagination() {
    const prevBtn = document.getElementById('prevBtn');
    const nextBtn = document.getElementById('nextBtn');

    if (prevBtn) {
        prevBtn.addEventListener('click', () => {
            if (currentPage > 0) {
                currentPage--;
                loadMyReports();
            }
        });
    }

    if (nextBtn) {
        nextBtn.addEventListener('click', () => {
            currentPage++;
            loadMyReports();
        });
    }
}

async function loadMyReports() {
    const grid = document.getElementById('myReportsGrid');
    if (!grid) return;

    grid.innerHTML = '<div class="loading">Loading your reports...</div>';

    try {
        const params = {
            limit: pageSize,
            offset: currentPage * pageSize
        };

        const data = await reportsApi.getUserReports(params);

        if (data.reports && data.reports.length > 0) {
            grid.innerHTML = '';
            data.reports.forEach(report => {
                grid.appendChild(createReportCard(report));
            });

            updatePagination(data.reports.length);
        } else {
            grid.innerHTML = `
                <div class="loading">
                    <p>You haven't created any reports yet.</p>
                    <a href="create-report.html" class="btn btn-primary" style="margin-top: 1rem;">
                        Create Your First Report
                    </a>
                </div>
            `;
            updatePagination(0);
        }
    } catch (error) {
        console.error('Error loading reports:', error);
        grid.innerHTML = '<div class="loading">Failed to load your reports</div>';
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
