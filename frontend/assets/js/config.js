// API Configuration
const API_CONFIG = {
    BASE_URL: 'https://ecommerce-api-go-production-1c31.up.railway.app',
    VERSION: 'v1',
    ENDPOINTS: {
        // Auth
        REGISTER: '/v1/user/register',
        LOGIN: '/v1/user/login',
        PROFILE: '/v1/user/me',
        ADMIN_PROFILE: '/v1/admin/me',
        
        // Reports
        REPORTS: '/v1/reports',
        REPORT_BY_ID: (id) => `/v1/reports/${id}`,
        USER_REPORTS: '/v1/user/reports',
        CREATE_REPORT: '/v1/reports',
        UPDATE_REPORT: (id) => `/v1/reports/${id}`,
        
        // Stats & Leaderboard
        STATS: '/v1/reports/stats',
        LEADERBOARD: '/v1/leaderboard',
        
        // Health
        HEALTH: '/v1/healthcheck'
    }
};

// Local Storage Keys
const STORAGE_KEYS = {
    AUTH_TOKEN: 'cityStars_authToken',
    USER_ROLE: 'cityStars_userRole',
    USER_ID: 'cityStars_userId'
};

// Category Icons
const CATEGORY_ICONS = {
    pothole: '🕳️',
    streetlight: '💡',
    water: '💧',
    garbage: '🗑️',
    road: '🛣️',
    other: '📝'
};

// Status Colors
const STATUS_CLASSES = {
    pending: 'pending',
    'in-progress': 'in-progress',
    completed: 'completed',
    rejected: 'rejected'
};
