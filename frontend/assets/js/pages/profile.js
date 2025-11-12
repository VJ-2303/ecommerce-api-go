// Profile Page

document.addEventListener('DOMContentLoaded', () => {
    // Require authentication
    if (!requireAuth()) return;

    setupMobileNav();
    loadProfile();
    loadUserStats();
});

async function loadProfile() {
    const profileInfo = document.getElementById('profileInfo');
    if (!profileInfo) return;

    try {
        const data = await authApi.getProfile();
        
        if (data) {
            const role = data.role || 'user';
            const userId = data.userID || data.AdminID || 'N/A';
            
            localStorage.setItem(STORAGE_KEYS.USER_ID, userId);
            localStorage.setItem(STORAGE_KEYS.USER_ROLE, role);

            profileInfo.innerHTML = `
                <h2>User #${userId}</h2>
                <div style="display: flex; gap: 0.5rem; align-items: center; margin-top: 0.5rem;">
                    <span class="report-status ${role === 'admin' ? 'completed' : 'in-progress'}" 
                          style="font-size: 0.875rem;">
                        ${role === 'admin' ? '👑 Admin' : '👤 User'}
                    </span>
                </div>
            `;
        }
    } catch (error) {
        console.error('Error loading profile:', error);
        profileInfo.innerHTML = '<div class="loading">Failed to load profile</div>';
    }
}

async function loadUserStats() {
    try {
        const data = await reportsApi.getUserReports({ limit: 100, offset: 0 });
        
        if (data.reports) {
            const reports = data.reports;
            const total = reports.length;
            const pending = reports.filter(r => r.status === 'pending').length;
            const completed = reports.filter(r => r.status === 'completed').length;

            document.getElementById('userTotalReports').textContent = total;
            document.getElementById('userPendingReports').textContent = pending;
            document.getElementById('userCompletedReports').textContent = completed;
        }
    } catch (error) {
        console.error('Error loading user stats:', error);
        // Keep default "-" values
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
