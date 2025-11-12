// Statistics Page

document.addEventListener('DOMContentLoaded', () => {
    setupMobileNav();
    loadStats();
});

async function loadStats() {
    try {
        const data = await statsApi.getStats();
        
        if (data.stats) {
            displayStats(data.stats);
            updateProgress(data.stats);
        }
    } catch (error) {
        console.error('Error loading stats:', error);
        showToast('Failed to load statistics', 'error');
    }
}

function displayStats(stats) {
    document.getElementById('totalReports').textContent = stats.total_reports || 0;
    document.getElementById('pendingReports').textContent = stats.pending_reports || 0;
    document.getElementById('inProgressReports').textContent = stats.in_progress_reports || 0;
    document.getElementById('completedReports').textContent = stats.completed_reports || 0;
    document.getElementById('rejectedReports').textContent = stats.rejected_reports || 0;
}

function updateProgress(stats) {
    const total = stats.total_reports || 0;
    const completed = stats.completed_reports || 0;
    
    let percentage = 0;
    if (total > 0) {
        percentage = Math.round((completed / total) * 100);
    }

    const progressFill = document.getElementById('progressFill');
    const progressPercent = document.getElementById('progressPercent');

    if (progressFill) {
        progressFill.style.width = `${percentage}%`;
    }

    if (progressPercent) {
        progressPercent.textContent = `${percentage}%`;
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
