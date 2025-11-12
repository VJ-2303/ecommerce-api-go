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

    // Load stats
    loadStats();

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
