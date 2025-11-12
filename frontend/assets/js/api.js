// API Service
class ApiService {
    constructor(baseUrl) {
        this.baseUrl = baseUrl;
    }

    async request(endpoint, options = {}) {
        const url = `${this.baseUrl}${endpoint}`;
        const token = localStorage.getItem(STORAGE_KEYS.AUTH_TOKEN);

        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };

        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const config = {
            ...options,
            headers
        };

        try {
            const response = await fetch(url, config);
            const data = await response.json();

            if (!response.ok) {
                throw data;
            }

            return data;
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    }

    // GET request
    async get(endpoint, params = {}) {
        const queryString = new URLSearchParams(params).toString();
        const url = queryString ? `${endpoint}?${queryString}` : endpoint;
        return this.request(url, { method: 'GET' });
    }

    // POST request
    async post(endpoint, data) {
        return this.request(endpoint, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }

    // PATCH request
    async patch(endpoint, data) {
        return this.request(endpoint, {
            method: 'PATCH',
            body: JSON.stringify(data)
        });
    }

    // DELETE request
    async delete(endpoint) {
        return this.request(endpoint, {
            method: 'DELETE'
        });
    }
}

// Initialize API service
const api = new ApiService(API_CONFIG.BASE_URL);

// Auth APIs
const authApi = {
    register: (data) => api.post(API_CONFIG.ENDPOINTS.REGISTER, data),
    login: (data) => api.post(API_CONFIG.ENDPOINTS.LOGIN, data),
    getProfile: () => api.get(API_CONFIG.ENDPOINTS.PROFILE),
    getAdminProfile: () => api.get(API_CONFIG.ENDPOINTS.ADMIN_PROFILE)
};

// Reports APIs
const reportsApi = {
    getAll: (params) => api.get(API_CONFIG.ENDPOINTS.REPORTS, params),
    getById: (id) => api.get(API_CONFIG.ENDPOINTS.REPORT_BY_ID(id)),
    getUserReports: (params) => api.get(API_CONFIG.ENDPOINTS.USER_REPORTS, params),
    create: (data) => api.post(API_CONFIG.ENDPOINTS.CREATE_REPORT, data),
    update: (id, data) => api.patch(API_CONFIG.ENDPOINTS.UPDATE_REPORT(id), data)
};

// Stats & Leaderboard APIs
const statsApi = {
    getStats: () => api.get(API_CONFIG.ENDPOINTS.STATS),
    getLeaderboard: () => api.get(API_CONFIG.ENDPOINTS.LEADERBOARD)
};

// Health API
const healthApi = {
    check: () => api.get(API_CONFIG.ENDPOINTS.HEALTH)
};
