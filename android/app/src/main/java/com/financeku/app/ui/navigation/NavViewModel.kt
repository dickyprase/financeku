package com.financeku.app.ui.navigation

import androidx.lifecycle.ViewModel
import com.financeku.app.data.local.datastore.TokenDataStore
import dagger.hilt.android.lifecycle.HiltViewModel
import javax.inject.Inject

@HiltViewModel
class NavViewModel @Inject constructor(
    val tokenDataStore: TokenDataStore
) : ViewModel()
