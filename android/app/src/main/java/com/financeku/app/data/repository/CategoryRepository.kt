package com.financeku.app.data.repository

import com.financeku.app.data.api.ApiService
import com.financeku.app.data.api.model.*
import com.financeku.app.data.local.dao.CategoryDao
import com.financeku.app.data.local.entity.CategoryEntity
import com.financeku.app.domain.model.Category
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class CategoryRepository @Inject constructor(
    private val apiService: ApiService,
    private val categoryDao: CategoryDao
) {
    fun getLocalCategories(): Flow<List<Category>> {
        return categoryDao.getAll().map { entities ->
            entities.map { it.toDomain() }
        }
    }

    fun getLocalCategoriesByType(type: String): Flow<List<Category>> {
        return categoryDao.getByType(type).map { entities ->
            entities.map { it.toDomain() }
        }
    }

    suspend fun fetchCategories(): Resource<List<Category>> {
        return try {
            val response = apiService.getCategories()
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data ?: emptyList()
                categoryDao.deleteAll()
                categoryDao.insertAll(data.map { it.toEntity() })
                Resource.Success(data.map { it.toDomain() })
            } else {
                Resource.Error(response.body()?.message ?: "Failed to fetch categories")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun createCategory(name: String, type: String, icon: String?, color: String?): Resource<Category> {
        return try {
            val response = apiService.createCategory(CategoryRequest(name, type, icon, color))
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                categoryDao.insert(data.toEntity())
                Resource.Success(data.toDomain())
            } else {
                Resource.Error(response.body()?.message ?: "Failed to create category")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun updateCategory(id: String, name: String, type: String, icon: String?, color: String?): Resource<Category> {
        return try {
            val response = apiService.updateCategory(id, CategoryRequest(name, type, icon, color))
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                categoryDao.insert(data.toEntity())
                Resource.Success(data.toDomain())
            } else {
                Resource.Error(response.body()?.message ?: "Failed to update category")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun deleteCategory(id: String): Resource<Unit> {
        return try {
            val response = apiService.deleteCategory(id)
            if (response.isSuccessful && response.body()?.success == true) {
                categoryDao.deleteById(id)
                Resource.Success(Unit)
            } else {
                Resource.Error(response.body()?.message ?: "Failed to delete category")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    private fun CategoryDto.toDomain() = Category(
        id = id, name = name, type = type, icon = icon, color = color
    )

    private fun CategoryDto.toEntity() = CategoryEntity(
        id = id, name = name, type = type, icon = icon, color = color
    )

    private fun CategoryEntity.toDomain() = Category(
        id = id, name = name, type = type, icon = icon, color = color
    )
}
