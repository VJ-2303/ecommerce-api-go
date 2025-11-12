// Main JavaScript for Homepage

// Mobile Navigation Toggle
document.addEventListener('DOMContentLoaded', () => {
    const navToggle = document.getElementById('navToggle');
    const navMenu = document.getElementById('navMenu');

    if (navToggle && navMenu) {
        navToggle.addEventListener('click', () => {
            navMenu.classList.toggle('active');
        });
    }

    // Close mobile menu when clicking outside
    document.addEventListener('click', (e) => {
        if (navMenu && navMenu.classList.contains('active')) {
            if (!navToggle.contains(e.target) && !navMenu.contains(e.target)) {
                navMenu.classList.remove('active');
            }
        }
    });

    // Load stats and recent reports
    loadStats();
    loadRecentReports();

    // Update hero button based on auth
    updateHeroButton();
});

// Update hero button based on authentication
function updateHeroButton() {
    const heroReportBtn = document.getElementById('heroReportBtn');
    if (!heroReportBtn) return;

    if (!isAuthenticated()) {
        heroReportBtn.href = 'pages/login.html';
    }
}

// Load statistics
async function loadStats() {
    try {
        const data = await statsApi.getStats();
        
        if (data.stats) {
            const stats = data.stats;
            
            document.getElementById('totalReports').textContent = stats.total_reports || 0;
            document.getElementById('pendingReports').textContent = stats.pending_reports || 0;
            document.getElementById('inProgressReports').textContent = stats.in_progress_reports || 0;
            document.getElementById('completedReports').textContent = stats.completed_reports || 0;
        }
    } catch (error) {
        console.error('Error loading stats:', error);
        // Keep the default "-" values
    }
}

// Load recent reports
async function loadRecentReports() {
    const grid = document.getElementById('recentReportsGrid');
    if (!grid) return;

    try {
        const data = await reportsApi.getAll({ limit: 6, offset: 0 });
        
        if (data.reports && data.reports.length > 0) {
            grid.innerHTML = '';
            data.reports.forEach(report => {
                grid.appendChild(createReportCard(report));
            });
        } else {
            grid.innerHTML = '<div class="loading">No reports available</div>';
        }
    } catch (error) {
        console.error('Error loading reports:', error);
        grid.innerHTML = '<div class="loading">Failed to load reports</div>';
    }
}

// Create report card
function createReportCard(report) {
    const card = document.createElement('div');
    card.className = 'report-card';
    card.onclick = () => {
        window.location.href = `pages/report-detail.html?id=${report.id}`;
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
