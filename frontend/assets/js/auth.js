// Authentication Management

// Check if user is authenticated
function isAuthenticated() {
    return !!localStorage.getItem(STORAGE_KEYS.AUTH_TOKEN);
}

// Get user role
function getUserRole() {
    return localStorage.getItem(STORAGE_KEYS.USER_ROLE) || 'user';
}

// Check if user is admin
function isAdmin() {
    return getUserRole() === 'admin';
}

// Save auth data
function saveAuthData(token, role = 'user') {
    localStorage.setItem(STORAGE_KEYS.AUTH_TOKEN, token);
    localStorage.setItem(STORAGE_KEYS.USER_ROLE, role);
}

// Clear auth data
function clearAuthData() {
    localStorage.removeItem(STORAGE_KEYS.AUTH_TOKEN);
    localStorage.removeItem(STORAGE_KEYS.USER_ROLE);
    localStorage.removeItem(STORAGE_KEYS.USER_ID);
}

// Logout
function logout() {
    clearAuthData();
    showToast('Logged out successfully', 'success');
    setTimeout(() => {
        window.location.href = '../index.html';
    }, 1000);
}

// Redirect to login if not authenticated
function requireAuth() {
    if (!isAuthenticated()) {
        showToast('Please login to continue', 'warning');
        setTimeout(() => {
            const currentPath = window.location.pathname;
            if (!currentPath.includes('login.html')) {
                window.location.href = currentPath.includes('/pages/') 
                    ? 'login.html' 
                    : 'pages/login.html';
            }
        }, 1000);
        return false;
    }
    return true;
}

// Require admin role
function requireAdmin() {
    if (!requireAuth()) return false;
    
    if (!isAdmin()) {
        showToast('Admin access required', 'error');
        setTimeout(() => {
            window.location.href = '../index.html';
        }, 1000);
        return false;
    }
    return true;
}

// Update navigation based on auth state
function updateNavigation() {
    const navAuth = document.getElementById('navAuth');
    const navUser = document.getElementById('navUser');
    const logoutBtn = document.getElementById('logoutBtn');

    if (!navAuth || !navUser) return;

    if (isAuthenticated()) {
        toggleElement(navAuth, false);
        toggleElement(navUser, true);

        if (logoutBtn) {
            logoutBtn.onclick = logout;
        }
    } else {
        toggleElement(navAuth, true);
        toggleElement(navUser, false);
    }
}

// Parse JWT token (without verification)
function parseJWT(token) {
    try {
        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));

        return JSON.parse(jsonPayload);
    } catch (error) {
        console.error('Error parsing JWT:', error);
        return null;
    }
}

// Check if token is expired
function isTokenExpired() {
    const token = localStorage.getItem(STORAGE_KEYS.AUTH_TOKEN);
    if (!token) return true;

    const payload = parseJWT(token);
    if (!payload || !payload.exp) return true;

    const expiryTime = payload.exp * 1000; // Convert to milliseconds
    return Date.now() >= expiryTime;
}

// Validate session
function validateSession() {
    if (isAuthenticated() && isTokenExpired()) {
        clearAuthData();
        showToast('Session expired. Please login again.', 'warning');
        return false;
    }
    return true;
}

// Initialize auth on page load
document.addEventListener('DOMContentLoaded', () => {
    validateSession();
    updateNavigation();
});
