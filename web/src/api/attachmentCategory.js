import service from '@/utils/request'
// Category list
export const getCategoryList = () => {
    return service({
        url: '/attachmentCategory/getCategoryList',
        method: 'get',
    })
}

// Add / edit category
export const addCategory = (data) => {
    return service({
        url: '/attachmentCategory/addCategory',
        method: 'post',
        data
    })
}

// Delete category
export const deleteCategory = (data) => {
    return service({
        url: '/attachmentCategory/deleteCategory',
        method: 'post',
        data
    })
}
